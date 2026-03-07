package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davideimola/worky/cmd/worky/templates"
	"github.com/spf13/cobra"
)

// NewInitCmd returns the `worky init <name>` command.
func NewInitCmd() *cobra.Command {
	var module string

	c := &cobra.Command{
		Use:   "init <name>",
		Short: "Scaffold a new workshop project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			slug := slugify(name)
			if module == "" {
				module = "github.com/example/" + slug
			}
			return initProject(name, module, ".")
		},
	}

	c.Flags().StringVar(&module, "module", "", "Go module path (default: github.com/example/<slug>)")
	return c
}

type projectData struct {
	Name   string
	Slug   string
	Module string
}

// slugify converts a display name to a kebab-case slug.
func slugify(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), " ", "-"))
}

// initProject scaffolds a new workshop project in <baseDir>/<slug>/.
func initProject(name, module, baseDir string) error {
	slug := slugify(name)
	dest := filepath.Join(baseDir, slug)

	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("directory %q already exists", dest)
	}

	data := projectData{
		Name:   name,
		Slug:   slug,
		Module: module,
	}

	type entry struct {
		tmpl string
		dest string
	}

	files := []entry{
		{"files/project/main.go.tmpl", "main.go"},
		{"files/project/go.mod.tmpl", "go.mod"},
		{"files/project/hugo.toml.tmpl", "hugo.toml"},
		{"files/project/Makefile.tmpl", "Makefile"},
		{"files/project/.gitignore.tmpl", ".gitignore"},
		{"files/project/README.md.tmpl", "README.md"},
		{"files/project/docs/_index.md.tmpl", "docs/_index.md"},
		{"files/project/docs/00-setup/_index.md.tmpl", "docs/00-setup/_index.md"},
		{"files/project/.github/workflows/release.yml.tmpl", ".github/workflows/release.yml"},
	}

	tmplFS := templates.FS()
	for _, f := range files {
		if err := scaffold(tmplFS, f.tmpl, filepath.Join(dest, f.dest), data); err != nil {
			return fmt.Errorf("scaffold %s: %w", f.dest, err)
		}
	}

	fmt.Printf("Workshop %q created in %s/\n\n", name, dest)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", slug)
	fmt.Printf("  hugo mod init %s\n", module)
	fmt.Println("  hugo mod get github.com/davideimola/worky")
	fmt.Println("  hugo mod get github.com/geekdocs/geekdoc")
	fmt.Println("  go mod tidy")
	fmt.Println("  make serve")
	return nil
}
