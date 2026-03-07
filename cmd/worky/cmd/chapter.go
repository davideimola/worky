package cmd

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/davideimola/worky/cmd/worky/templates"
	"github.com/spf13/cobra"
)

// NewChapterCmd returns the `worky new chapter <id> <name>` command.
func NewChapterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chapter <id> <name>",
		Short: "Add a new chapter to the workshop",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newChapter(args[0], args[1], ".", cmd.OutOrStdout())
		},
	}
}

type chapterData struct {
	ID   string
	Name string
	Slug string
}

// newChapter scaffolds a new chapter under <baseDir>/site/<id>-<slug>/index.md
// and prints the Go snippet to add to main.go.
func newChapter(id, name, baseDir string, out io.Writer) error {
	slug := slugify(name)
	dirName := id + "-" + slug
	destPath := filepath.Join(baseDir, "site", dirName, "index.md")

	data := chapterData{
		ID:   id,
		Name: name,
		Slug: dirName,
	}

	if err := scaffold(templates.FS(), "files/chapter/index.md.tmpl", destPath, data); err != nil {
		return fmt.Errorf("scaffold chapter: %w", err)
	}

	_, _ = fmt.Fprintf(out, "Chapter created: %s\n\n", destPath)
	_, _ = fmt.Fprintf(out, "Add the following to your main.go Chapters slice:\n\n")
	_, _ = fmt.Fprintf(out, "    worky.Chapter{\n")
	_, _ = fmt.Fprintf(out, "        ID:   %q,\n", id)
	_, _ = fmt.Fprintf(out, "        Name: %q,\n", name)
	_, _ = fmt.Fprintf(out, "        Slug: %q,\n", dirName)
	_, _ = fmt.Fprintf(out, "        Checks: []worky.Check{\n")
	_, _ = fmt.Fprintf(out, "            // TODO: add checks here\n")
	_, _ = fmt.Fprintf(out, "        },\n")
	_, _ = fmt.Fprintf(out, "    },\n")

	return nil
}
