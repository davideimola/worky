package worky

import "testing"

func newTestWorkshop() *Workshop {
	return New(Config{
		Name: "test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
			{ID: "02", Name: "Two", Slug: "02-two"},
			{ID: "03", Name: "Three", Slug: "03-three"},
		},
	})
}

func TestChapterByID_Found(t *testing.T) {
	w := newTestWorkshop()
	c, ok := w.chapterByID("02")
	if !ok {
		t.Fatal("expected to find chapter 02")
	}
	if c.Name != "Two" {
		t.Errorf("expected name Two, got %q", c.Name)
	}
}

func TestChapterByID_NotFound(t *testing.T) {
	w := newTestWorkshop()
	_, ok := w.chapterByID("99")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestChapterBySlug_Found(t *testing.T) {
	w := newTestWorkshop()
	c, ok := w.chapterBySlug("02-two")
	if !ok {
		t.Fatal("expected to find chapter by slug")
	}
	if c.ID != "02" {
		t.Errorf("expected ID 02, got %q", c.ID)
	}
}

func TestChapterBySlug_NotFound(t *testing.T) {
	w := newTestWorkshop()
	_, ok := w.chapterBySlug("99-nope")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestNextChapter(t *testing.T) {
	w := newTestWorkshop()
	next, ok := w.nextChapter("01")
	if !ok {
		t.Fatal("expected next chapter after 01")
	}
	if next.ID != "02" {
		t.Errorf("expected 02, got %q", next.ID)
	}
}

func TestNextChapter_Last(t *testing.T) {
	w := newTestWorkshop()
	_, ok := w.nextChapter("03")
	if ok {
		t.Fatal("last chapter should have no next")
	}
}

func TestNextChapter_NotFound(t *testing.T) {
	w := newTestWorkshop()
	_, ok := w.nextChapter("99")
	if ok {
		t.Fatal("unknown chapter should have no next")
	}
}

func TestNew_DuplicateID_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for duplicate chapter ID")
		}
	}()
	New(Config{
		Name: "test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
			{ID: "01", Name: "Dupe", Slug: "01-dupe"},
		},
	})
}

func TestNew_DuplicateSlug_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for duplicate chapter slug")
		}
	}()
	New(Config{
		Name: "test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "same-slug"},
			{ID: "02", Name: "Two", Slug: "same-slug"},
		},
	})
}
