package worky

// Chapter represents a workshop chapter.
type Chapter struct {
	ID     string  // "00", "01", etc.
	Name   string  // Display name
	Slug   string  // URL slug, e.g. "00-setup"
	Checks []Check // Validation steps for this chapter
}

func (w *Workshop) chapterByID(id string) (Chapter, bool) {
	i, ok := w.idxByID[id]
	if !ok {
		return Chapter{}, false
	}
	return w.cfg.Chapters[i], true
}

func (w *Workshop) chapterBySlug(slug string) (Chapter, bool) {
	i, ok := w.idxBySlug[slug]
	if !ok {
		return Chapter{}, false
	}
	return w.cfg.Chapters[i], true
}

func (w *Workshop) nextChapter(id string) (Chapter, bool) {
	i, ok := w.idxByID[id]
	if !ok || i+1 >= len(w.cfg.Chapters) {
		return Chapter{}, false
	}
	return w.cfg.Chapters[i+1], true
}
