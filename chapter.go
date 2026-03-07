package worky

// Chapter represents a workshop chapter.
type Chapter struct {
	ID     string  // "00", "01", etc.
	Name   string  // Display name
	Slug   string  // URL slug, e.g. "00-setup"
	Checks []Check // Validation steps for this chapter
}

func (w *Workshop) chapterByID(id string) (Chapter, bool) {
	for _, c := range w.cfg.Chapters {
		if c.ID == id {
			return c, true
		}
	}
	return Chapter{}, false
}

func (w *Workshop) chapterBySlug(slug string) (Chapter, bool) {
	for _, c := range w.cfg.Chapters {
		if c.Slug == slug {
			return c, true
		}
	}
	return Chapter{}, false
}

func (w *Workshop) nextChapter(id string) (Chapter, bool) {
	for i, c := range w.cfg.Chapters {
		if c.ID == id && i+1 < len(w.cfg.Chapters) {
			return w.cfg.Chapters[i+1], true
		}
	}
	return Chapter{}, false
}
