package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"

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
	ID     string
	Name   string
	Slug   string
	Weight int
}

// newChapter scaffolds a new chapter under <baseDir>/docs/<id>-<slug>/_index.md
// and prints the Go snippet to add to main.go.
func newChapter(id, name, baseDir string, out io.Writer) error {
	slug := slugify(name)
	dirName := id + "-" + slug
	destPath := filepath.Join(baseDir, "docs", dirName, "_index.md")

	weight, _ := strconv.Atoi(id)

	data := chapterData{
		ID:     id,
		Name:   name,
		Slug:   slug,
		Weight: weight,
	}

	if err := scaffold(templates.FS(), "files/chapter/_index.md.tmpl", destPath, data); err != nil {
		return fmt.Errorf("scaffold chapter: %w", err)
	}

	fmt.Fprintf(out, "Chapter created: %s\n\n", destPath)
	fmt.Fprintf(out, "Add the following to your main.go Chapters slice:\n\n")
	fmt.Fprintf(out, "    worky.Chapter{\n")
	fmt.Fprintf(out, "        ID:   %q,\n", id)
	fmt.Fprintf(out, "        Name: %q,\n", name)
	fmt.Fprintf(out, "        Slug: %q,\n", dirName)
	fmt.Fprintf(out, "        Checks: []worky.Check{\n")
	fmt.Fprintf(out, "            // TODO: add checks here\n")
	fmt.Fprintf(out, "        },\n")
	fmt.Fprintf(out, "    },\n")

	return nil
}
