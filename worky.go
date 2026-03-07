package worky

import (
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Config holds the workshop configuration provided by the creator.
type Config struct {
	Name     string            // Workshop display name
	HomeDir  string            // State directory under home (default: ".worky")
	Port     int               // Default HTTP port (default: 8080)
	SiteFS   fs.FS             // Embedded Hugo output (must contain a "site/" subdirectory)
	Chapters []Chapter
}

// Workshop is the runnable workshop instance.
type Workshop struct {
	cfg Config
	hub *sseHub
}

// New creates a new Workshop with the given config, applying defaults.
func New(cfg Config) *Workshop {
	if cfg.HomeDir == "" {
		cfg.HomeDir = ".worky"
	}
	if cfg.Port == 0 {
		cfg.Port = 8080
	}
	return &Workshop{cfg: cfg, hub: newSSEHub()}
}

// Run builds the CLI and executes it.
func (w *Workshop) Run() {
	root := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: w.cfg.Name,
	}

	root.AddCommand(
		w.serveCmd(),
		w.checkCmd(),
		w.statusCmd(),
		w.resetCmd(),
		w.unlockCmd(),
		w.stopCmd(),
		w.logsCmd(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func (w *Workshop) serveCmd() *cobra.Command {
	var port int
	var open bool
	var detach bool
	var preview bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the workshop web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if detach {
				return w.runDetached(port, open)
			}

			if preview {
				fmt.Println("Preview mode: chapter locking is disabled.")
			}

			handler := w.newHandler(preview)
			addr := net.JoinHostPort("", strconv.Itoa(port))
			url := fmt.Sprintf("http://localhost:%d", port)

			fmt.Printf("Workshop server starting on %s\n", url)

			go w.watchFiles(cmd.Context())

			if open {
				go func() {
					time.Sleep(500 * time.Millisecond)
					openBrowser(url)
				}()
			}

			return http.ListenAndServe(addr, handler)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", w.cfg.Port, "Port to listen on")
	cmd.Flags().BoolVar(&open, "open", false, "Open browser automatically")
	cmd.Flags().BoolVar(&detach, "detach", false, "Run server in background")
	cmd.Flags().BoolVar(&preview, "preview", false, "Disable chapter locking (for content review)")
	return cmd
}

func (w *Workshop) runDetached(port int, open bool) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	if err == nil {
		conn.Close()
		fmt.Printf("Workshop server is already running on http://localhost:%d\n", port)
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	dir, err := w.workshopDir()
	if err != nil {
		return fmt.Errorf("failed to get workshop dir: %w", err)
	}

	logPath := filepath.Join(dir, "server.log")
	pidPath := filepath.Join(dir, "server.pid")

	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	cmdArgs := []string{"serve", "--port", strconv.Itoa(port)}
	if open {
		cmdArgs = append(cmdArgs, "--open")
	}

	c := exec.Command(exe, cmdArgs...)
	c.Stdout = logFile
	c.Stderr = logFile
	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(c.Process.Pid)), 0o644); err != nil {
		c.Process.Kill() //nolint:errcheck
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	fmt.Printf("Workshop server started (PID %d)\n", c.Process.Pid)
	fmt.Printf("Logs: %s\n", logPath)
	return nil
}

func (w *Workshop) stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the background workshop server",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := w.workshopDir()
			if err != nil {
				return fmt.Errorf("failed to get workshop dir: %w", err)
			}

			pidPath := filepath.Join(dir, "server.pid")
			data, err := os.ReadFile(pidPath)
			if os.IsNotExist(err) {
				fmt.Println("No server PID file found. Is the server running?")
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to read PID file: %w", err)
			}

			pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
			if err != nil {
				return fmt.Errorf("invalid PID in file: %w", err)
			}

			proc, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Errorf("failed to find process %d: %w", pid, err)
			}

			if err := proc.Kill(); err != nil {
				return fmt.Errorf("failed to kill process %d: %w", pid, err)
			}

			os.Remove(pidPath) //nolint:errcheck
			fmt.Printf("Workshop server (PID %d) stopped.\n", pid)
			return nil
		},
	}
}

func (w *Workshop) checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check [chapter]",
		Short: "Validate a chapter and unlock the next one",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := ""
			if len(args) > 0 {
				id = args[0]
			} else {
				state, _ := w.loadProgress()
				for _, c := range w.cfg.Chapters {
					if state.IsUnlocked(c.ID) && !state.IsCompleted(c.ID) {
						id = c.ID
						break
					}
				}
				if id == "" {
					fmt.Println("No chapter to check. All chapters are completed or none are unlocked.")
					return nil
				}
				fmt.Printf("Auto-detected current chapter: %s\n", id)
			}

			chapter, ok := w.chapterByID(id)
			if !ok {
				return fmt.Errorf("unknown chapter: %q", id)
			}

			if len(chapter.Checks) == 0 {
				return fmt.Errorf("no checks defined for chapter %s", id)
			}

			fmt.Printf("\nChecking Chapter %s: %s\n\n", chapter.ID, chapter.Name)
			results, passed := w.runChecks(chapter.Checks)

			if err := w.saveCheckResults(chapter.ID, results); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to save check results: %v\n", err)
			}

			for _, r := range results {
				if r.Passed {
					fmt.Printf("  \033[32m✓\033[0m %s\n", r.Description)
				} else {
					fmt.Printf("  \033[31m✗\033[0m %s\n    → %s\n", r.Description, r.Error)
				}
			}

			fmt.Println()

			if !passed {
				fmt.Println("\033[31mSome checks failed. Fix the issues above and try again.\033[0m")
				os.Exit(1)
			}

			state, err := w.loadProgress()
			if err != nil {
				return fmt.Errorf("failed to load progress: %w", err)
			}

			next, hasNext := w.nextChapter(id)
			nextID := ""
			if hasNext {
				nextID = next.ID
			}
			state.Complete(id, nextID)

			if err := w.saveProgress(state); err != nil {
				return fmt.Errorf("failed to save progress: %w", err)
			}

			fmt.Printf("\033[32mChapter %s complete!\033[0m\n", id)
			if hasNext {
				fmt.Printf("Chapter %s (%s) is now unlocked.\n", next.ID, next.Name)
				fmt.Printf("Open \033[34mhttp://localhost:%d/%s/\033[0m\n", w.cfg.Port, next.Slug)
			} else {
				fmt.Println("You have completed the entire workshop!")
			}

			return nil
		},
	}
}

func (w *Workshop) resetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset all progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := w.resetProgress(); err != nil {
				return fmt.Errorf("failed to reset progress: %w", err)
			}
			firstID := ""
			if len(w.cfg.Chapters) > 0 {
				firstID = w.cfg.Chapters[0].ID
			}
			fmt.Printf("Progress reset. Chapter %s is unlocked.\n", firstID)
			return nil
		},
	}
}

func (w *Workshop) unlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlock <chapter>",
		Short: "Manually unlock a chapter (facilitator use)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			chapter, ok := w.chapterByID(id)
			if !ok {
				return fmt.Errorf("unknown chapter: %q", id)
			}

			state, err := w.loadProgress()
			if err != nil {
				return fmt.Errorf("failed to load progress: %w", err)
			}

			if state.IsUnlocked(id) {
				fmt.Printf("Chapter %s (%s) is already unlocked.\n", chapter.ID, chapter.Name)
				return nil
			}

			state.Unlock(id)
			if err := w.saveProgress(state); err != nil {
				return fmt.Errorf("failed to save progress: %w", err)
			}

			fmt.Printf("Chapter %s (%s) unlocked.\n", chapter.ID, chapter.Name)
			fmt.Printf("Open \033[34mhttp://localhost:%d/%s/\033[0m\n", w.cfg.Port, chapter.Slug)
			return nil
		},
	}
}

func (w *Workshop) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the status of all chapters",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := w.loadProgress()
			if err != nil {
				return fmt.Errorf("failed to load progress: %w", err)
			}

			fmt.Printf("\n%s — Chapter Status\n\n", w.cfg.Name)
			for _, c := range w.cfg.Chapters {
				icon := "🔒"
				if state.IsCompleted(c.ID) {
					icon = "✅"
				} else if state.IsUnlocked(c.ID) {
					icon = "🔓"
				}
				fmt.Printf("  %s  Chapter %s: %s\n", icon, c.ID, c.Name)
			}
			fmt.Println()
			return nil
		},
	}
}

func (w *Workshop) logsCmd() *cobra.Command {
	var follow bool
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show workshop server logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := w.workshopDir()
			if err != nil {
				return fmt.Errorf("failed to get workshop dir: %w", err)
			}

			logPath := filepath.Join(dir, "server.log")
			f, err := os.Open(logPath)
			if os.IsNotExist(err) {
				fmt.Println("No log file found. Has the server been started with --detach?")
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to open log file: %w", err)
			}
			defer f.Close()

			io.Copy(os.Stdout, f) //nolint:errcheck
			if !follow {
				return nil
			}
			for {
				time.Sleep(500 * time.Millisecond)
				io.Copy(os.Stdout, f) //nolint:errcheck
			}
		},
	}
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	return cmd
}

func openBrowser(url string) {
	var cmd string
	var cmdArgs []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		cmdArgs = []string{url}
	case "windows":
		cmd = "cmd"
		cmdArgs = []string{"/c", "start", url}
	default:
		cmd = "xdg-open"
		cmdArgs = []string{url}
	}
	exec.Command(cmd, cmdArgs...).Start() //nolint:errcheck
}
