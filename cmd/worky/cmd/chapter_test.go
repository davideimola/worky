package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewChapter_CreatesMarkdownFile(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	if err := newChapter("01", "Hello World", dir, &out); err != nil {
		t.Fatalf("newChapter: %v", err)
	}

	path := filepath.Join(dir, "docs", "01-hello-world", "_index.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file %s to exist: %v", path, err)
	}
}

func TestNewChapter_SlugFromName(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	if err := newChapter("02", "Hello World", dir, &out); err != nil {
		t.Fatalf("newChapter: %v", err)
	}

	path := filepath.Join(dir, "docs", "02-hello-world", "_index.md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected directory with slug 'hello-world': %v", err)
	}
}

func TestNewChapter_FileContents(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	if err := newChapter("03", "My Chapter", dir, &out); err != nil {
		t.Fatalf("newChapter: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "docs", "03-my-chapter", "_index.md"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), `title: "My Chapter"`) {
		t.Errorf("expected chapter name in file, got:\n%s", content)
	}
}

func TestNewChapter_PrintsGoSnippet(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer

	if err := newChapter("04", "My Chapter", dir, &out); err != nil {
		t.Fatalf("newChapter: %v", err)
	}

	snippet := out.String()
	if !strings.Contains(snippet, `ID:   "04"`) {
		t.Errorf("expected ID in snippet, got:\n%s", snippet)
	}
	if !strings.Contains(snippet, `Name: "My Chapter"`) {
		t.Errorf("expected Name in snippet, got:\n%s", snippet)
	}
	if !strings.Contains(snippet, `Slug: "04-my-chapter"`) {
		t.Errorf("expected Slug in snippet, got:\n%s", snippet)
	}
}
