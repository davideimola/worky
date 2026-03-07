package worky

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

// sseHub broadcasts update signals to all connected SSE clients.
type sseHub struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newSSEHub() *sseHub {
	return &sseHub{clients: make(map[chan struct{}]struct{})}
}

func (h *sseHub) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *sseHub) unsubscribe(ch chan struct{}) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
}

func (h *sseHub) broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// watchFiles polls progress.json and checks.json every 500ms and broadcasts
// on any change, so connected SSE clients receive instant updates.
func (w *Workshop) watchFiles(ctx context.Context) {
	var lastMod time.Time
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mod := w.latestFileMod()
			if !mod.IsZero() && mod.After(lastMod) {
				lastMod = mod
				w.hub.broadcast()
			}
		}
	}
}

func (w *Workshop) latestFileMod() time.Time {
	var latest time.Time
	for _, getter := range []func() (string, error){w.stateFile, w.checkResultsFile} {
		path, err := getter()
		if err != nil {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
	}
	return latest
}

const lockedPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Chapter Locked</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
           display: flex; align-items: center; justify-content: center;
           min-height: 100vh; margin: 0; background: #1e1e2e; color: #cdd6f4; }
    .card { text-align: center; padding: 3rem 4rem; background: #313244;
            border-radius: 1rem; box-shadow: 0 8px 32px rgba(0,0,0,0.4); max-width: 480px; }
    .lock { font-size: 4rem; margin-bottom: 1rem; }
    h1 { margin: 0 0 0.5rem; font-size: 1.8rem; }
    p { color: #a6adc8; margin: 0.5rem 0; }
    code { background: #45475a; padding: 0.2em 0.5em; border-radius: 0.3em;
           font-family: monospace; color: #f5c2e7; }
  </style>
</head>
<body>
  <div class="card">
    <div class="lock">🔒</div>
    <h1>Chapter Locked</h1>
    <p>Complete the previous chapter first.</p>
    <p>Run <code>%s</code> to verify your work and unlock this chapter.</p>
    <p style="margin-top:1.5rem;font-size:0.85rem">
      <a href="/" style="color:#89b4fa">← Back to Home</a>
    </p>
  </div>
</body>
</html>`

// ChapterStatus is used in the API response.
type ChapterStatus struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Unlocked  bool   `json:"unlocked"`
	Completed bool   `json:"completed"`
}

// ProgressResponse is the /api/progress response.
type ProgressResponse struct {
	Completed []string        `json:"completed"`
	Unlocked  []string        `json:"unlocked"`
	Chapters  []ChapterStatus `json:"chapters"`
}

func (w *Workshop) newHandler(preview bool) http.Handler {
	var subFS fs.FS
	var fileServer http.Handler
	if w.cfg.SiteFS != nil {
		sub, err := fs.Sub(w.cfg.SiteFS, "site")
		if err != nil {
			panic(fmt.Sprintf("failed to sub site FS: %v", err))
		}
		subFS = sub
		fileServer = http.FileServer(http.FS(subFS))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/progress", w.handleProgress)
	mux.HandleFunc("/api/checks", w.handleChecks)
	mux.HandleFunc("/api/events", w.handleEvents)
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if !preview {
			if chapter, locked := w.isLocked(r.URL.Path); locked {
				prevID := ""
				for i, c := range w.cfg.Chapters {
					if c.Slug == chapter.Slug && i > 0 {
						prevID = w.cfg.Chapters[i-1].ID
						break
					}
				}
				cmd := os.Args[0] + " check " + prevID
				rw.WriteHeader(http.StatusForbidden)
				_, _ = fmt.Fprintf(rw, lockedPage, cmd)
				return
			}
		}
		if w.serveMD(rw, r, subFS, preview) {
			return
		}
		if fileServer != nil {
			fileServer.ServeHTTP(rw, r)
			return
		}
		w.serveBuiltinHome(rw, r)
	})

	return mux
}

func (w *Workshop) isLocked(urlPath string) (Chapter, bool) {
	path := strings.TrimPrefix(urlPath, "/")
	seg := strings.SplitN(path, "/", 2)[0]
	if seg == "" {
		return Chapter{}, false
	}

	chapter, ok := w.chapterBySlug(seg)
	if !ok {
		return Chapter{}, false
	}

	state, err := w.loadProgress()
	if err != nil {
		return Chapter{}, false
	}

	if !state.IsUnlocked(chapter.ID) {
		return chapter, true
	}
	return Chapter{}, false
}

const mdPageTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.ChapterID}} — {{.ChapterName}} | {{.WorkshopName}}</title>
  <link rel="stylesheet" href="/workshop.css">
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #1e1e2e; color: #cdd6f4; display: flex; min-height: 100vh; }
    nav { width: 260px; min-height: 100vh; background: #181825; border-right: 1px solid #313244; display: flex; flex-direction: column; flex-shrink: 0; }
    .nav-logo { padding: 1.25rem 1.5rem; font-size: 1.05rem; font-weight: 700; color: #cba6f7; border-bottom: 1px solid #313244; text-decoration: none; }
    .nav-section { padding: 1rem 0; }
    .nav-label { padding: 0.35rem 1.5rem; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.08em; color: #6c7086; }
    a.gdoc-nav__entry { display: block; padding: 0.45rem 1.5rem; color: #bac2de; text-decoration: none; font-size: 0.9rem; border-left: 3px solid transparent; transition: color 0.15s, background 0.15s; }
    a.gdoc-nav__entry:hover { color: #cdd6f4; background: #313244; }
    a.gdoc-nav__entry.active { color: #cba6f7; border-left-color: #cba6f7; background: rgba(203,166,247,0.08); }
    main { flex: 1; padding: 3rem 4rem; max-width: 860px; }
    h1 { font-size: 1.8rem; color: #cdd6f4; margin-bottom: 1.5rem; padding-bottom: 1.5rem; border-bottom: 1px solid #313244; }
    h2 { font-size: 1.15rem; color: #89b4fa; margin: 2rem 0 0.75rem; }
    h3 { font-size: 1rem; color: #89b4fa; margin: 1.5rem 0 0.5rem; }
    p { line-height: 1.7; color: #bac2de; margin-bottom: 1rem; }
    a { color: #89b4fa; }
    code { background: #313244; color: #f5c2e7; padding: 0.15em 0.45em; border-radius: 4px; font-family: "JetBrains Mono","Fira Code",monospace; font-size: 0.88em; }
    pre { background: #181825; border: 1px solid #313244; border-radius: 8px; padding: 1rem 1.25rem; overflow-x: auto; margin-bottom: 1rem; }
    pre code { background: none; padding: 0; color: #cdd6f4; font-size: 0.87rem; }
    ul, ol { color: #bac2de; line-height: 1.7; margin-bottom: 1rem; padding-left: 1.5rem; }
    li { margin-bottom: 0.25rem; }
    hr { border: none; border-top: 1px solid #313244; margin: 1.5rem 0; }
    blockquote { border-left: 3px solid #89b4fa; padding: 0.5rem 1rem; color: #a6adc8; margin-bottom: 1rem; background: rgba(137,180,250,0.06); border-radius: 0 4px 4px 0; }
    blockquote p { margin: 0; }
    details { border-radius: 6px; margin: 1rem 0; overflow: hidden; }
    details[data-type="hint"] { border: 1px solid #3b82f6; background: rgba(59,130,246,0.08); }
    details[data-type="solution"] { border: 1px solid #22c55e; background: rgba(34,197,94,0.08); }
    details summary { list-style: none; padding: 0.65rem 1rem; cursor: pointer; font-weight: 600; user-select: none; display: flex; align-items: center; gap: 0.5rem; }
    details summary::-webkit-details-marker { display: none; }
    details[data-type="hint"] summary { color: #93c5fd; }
    details[data-type="solution"] summary { color: #86efac; }
    details .details-body { padding: 0.25rem 1rem 0.75rem; border-top: 1px solid rgba(255,255,255,0.06); }
    details .details-body p:last-child { margin-bottom: 0; }
  </style>
</head>
<body>
  {{if .Preview}}<div style="position:sticky;top:0;z-index:1000;background:#fab387;color:#1e1e2e;padding:0.4rem 1.5rem;font-size:0.82rem;font-weight:600;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;display:flex;align-items:center;gap:0.5rem;">
    <span>&#9711;</span> Preview mode — chapter locking disabled
  </div>{{end}}
  <nav>
    <a class="nav-logo" href="/">⚙ {{.WorkshopName}}</a>
    <div class="nav-section">
      <div class="nav-label">Chapters</div>
      {{range .Chapters}}<a class="gdoc-nav__entry{{if eq .Slug $.CurrentSlug}} active{{end}}" href="/{{.Slug}}/">{{.ID}} — {{.Name}}</a>
      {{end}}
    </div>
  </nav>
  <main>
    {{.Content}}
    <div id="ws-check-results" data-chapter-id="{{.ChapterID}}" style="display:none"></div>
  </main>
  <script src="/workshop-progress.js"></script>
</body>
</html>`

type mdPageData struct {
	WorkshopName string
	ChapterID    string
	ChapterName  string
	CurrentSlug  string
	Chapters     []Chapter
	Content      template.HTML
	Preview      bool
}

var mdTmpl = template.Must(template.New("md").Parse(mdPageTmpl))

func (w *Workshop) serveMD(rw http.ResponseWriter, r *http.Request, subFS fs.FS, preview bool) bool {
	if subFS == nil {
		return false
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	seg := strings.SplitN(path, "/", 2)[0]
	if seg == "" {
		return false
	}

	chapter, ok := w.chapterBySlug(seg)
	if !ok {
		return false
	}

	mdData, err := fs.ReadFile(subFS, seg+"/index.md")
	if err != nil {
		return false
	}

	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
	if err := md.Convert(mdData, &buf); err != nil {
		return false
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mdTmpl.Execute(rw, mdPageData{
		WorkshopName: w.cfg.Name,
		ChapterID:    chapter.ID,
		ChapterName:  chapter.Name,
		CurrentSlug:  chapter.Slug,
		Chapters:     w.cfg.Chapters,
		Content:      template.HTML(buf.String()),
		Preview:      preview,
	}); err != nil {
		log.Printf("worky: failed to render chapter page for %q: %v", seg, err)
	}
	return true
}

const builtinHomeTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.WorkshopName}}</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
           background: #1e1e2e; color: #cdd6f4; min-height: 100vh;
           display: flex; align-items: flex-start; justify-content: center; padding: 4rem 1rem; }
    .container { width: 100%; max-width: 640px; }
    h1 { font-size: 2rem; color: #cba6f7; margin-bottom: 0.5rem; }
    .subtitle { color: #a6adc8; font-size: 0.95rem; margin-bottom: 2.5rem; }
    .chapter-list { list-style: none; display: flex; flex-direction: column; gap: 0.75rem; }
    .chapter-item { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.25rem;
                    background: #313244; border-radius: 0.75rem; border: 1px solid #45475a;
                    transition: border-color 0.15s; }
    .chapter-item.unlocked { border-color: #89b4fa; }
    .chapter-item.completed { border-color: #a6e3a1; opacity: 0.8; }
    .chapter-item.locked .chapter-name { color: #6c7086; }
    .chapter-icon { font-size: 1.4rem; flex-shrink: 0; }
    .chapter-id { font-size: 0.75rem; font-weight: 700; text-transform: uppercase;
                  letter-spacing: 0.06em; color: #6c7086; margin-bottom: 0.2rem; }
    .chapter-name { font-size: 1rem; color: #cdd6f4; }
  </style>
</head>
<body>
  <div class="container">
    <h1>⚙ {{.WorkshopName}}</h1>
    <p class="subtitle">Workshop progress</p>
    <ul class="chapter-list" id="chapter-list">
      {{range .Chapters}}
      <li class="chapter-item{{if .Completed}} completed{{else if .Unlocked}} unlocked{{else}} locked{{end}}" data-chapter-id="{{.ID}}">
        <span class="chapter-icon">{{if .Completed}}✅{{else if .Unlocked}}🔓{{else}}🔒{{end}}</span>
        <div>
          <div class="chapter-id">Chapter {{.ID}}</div>
          <div class="chapter-name">{{.Name}}</div>
        </div>
      </li>
      {{end}}
    </ul>
  </div>
  <script>
    const es = new EventSource('/api/events');
    es.onmessage = function() {
      fetch('/api/progress').then(r => r.json()).then(data => {
        data.chapters.forEach(ch => {
          const el = document.querySelector('[data-chapter-id="' + ch.id + '"]');
          if (!el) return;
          el.className = 'chapter-item ' + (ch.completed ? 'completed' : ch.unlocked ? 'unlocked' : 'locked');
          el.querySelector('.chapter-icon').textContent = ch.completed ? '✅' : ch.unlocked ? '🔓' : '🔒';
        });
      });
    };
  </script>
</body>
</html>`

type builtinHomeData struct {
	WorkshopName string
	Chapters     []ChapterStatus
}

var builtinHomeTpl = template.Must(template.New("home").Parse(builtinHomeTmpl))

func (w *Workshop) serveBuiltinHome(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}
	state, _ := w.loadProgress()
	var chapters []ChapterStatus
	for _, c := range w.cfg.Chapters {
		cs := ChapterStatus{ID: c.ID, Name: c.Name, Slug: c.Slug}
		if state != nil {
			cs.Unlocked = state.IsUnlocked(c.ID)
			cs.Completed = state.IsCompleted(c.ID)
		}
		chapters = append(chapters, cs)
	}
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := builtinHomeTpl.Execute(rw, builtinHomeData{
		WorkshopName: w.cfg.Name,
		Chapters:     chapters,
	}); err != nil {
		log.Printf("worky: failed to render home page: %v", err)
	}
}

func (w *Workshop) handleEvents(rw http.ResponseWriter, r *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "SSE not supported", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	ch := w.hub.subscribe()
	defer w.hub.unsubscribe(ch)

	_, _ = fmt.Fprintf(rw, "data: connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			_, _ = fmt.Fprintf(rw, "data: update\n\n")
			flusher.Flush()
		}
	}
}

func (w *Workshop) handleChecks(rw http.ResponseWriter, r *http.Request) {
	store, err := w.loadCheckResults()
	if err != nil {
		http.Error(rw, "failed to load check results", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(store); err != nil {
		log.Printf("worky: failed to encode check results: %v", err)
	}
}

func (w *Workshop) handleProgress(rw http.ResponseWriter, r *http.Request) {
	state, err := w.loadProgress()
	if err != nil {
		http.Error(rw, "failed to load progress", http.StatusInternalServerError)
		return
	}

	resp := ProgressResponse{
		Completed: state.Completed,
		Unlocked:  state.Unlocked,
	}
	if resp.Completed == nil {
		resp.Completed = []string{}
	}
	if resp.Unlocked == nil {
		resp.Unlocked = []string{}
	}

	for _, c := range w.cfg.Chapters {
		resp.Chapters = append(resp.Chapters, ChapterStatus{
			ID:        c.ID,
			Name:      c.Name,
			Slug:      c.Slug,
			Unlocked:  state.IsUnlocked(c.ID),
			Completed: state.IsCompleted(c.ID),
		})
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		log.Printf("worky: failed to encode progress response: %v", err)
	}
}
