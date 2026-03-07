package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davideimola/worky/cmd/worky/templates"
	"github.com/spf13/cobra"
)

// NewInitCmd returns the `worky init [name]` command.
func NewInitCmd() *cobra.Command {
	var (
		module        string
		yes           bool
		noInstallDeps bool
	)

	c := &cobra.Command{
		Use:   "init [name]",
		Short: "Scaffold a new workshop project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var prompter Prompter
			if yes {
				prompter = &yesPrompter{out: cmd.OutOrStdout()}
			} else {
				prompter = newStdinPrompter(cmd.InOrStdin(), cmd.OutOrStdout())
			}

			// Resolve name.
			var name string
			if len(args) > 0 {
				name = args[0]
			} else {
				var err error
				name, err = prompter.Ask("Workshop name", "")
				if err != nil {
					return err
				}
				if name == "" {
					return errors.New("workshop name is required")
				}
			}

			slug := slugify(name)

			// Resolve module path.
			defaultModule := "github.com/example/" + slug
			if !cmd.Flags().Changed("module") {
				var err error
				module, err = prompter.Ask("Go module path", defaultModule)
				if err != nil {
					return err
				}
				if module == "" {
					module = defaultModule
				}
			}

			// Resolve installDeps.
			installDeps := false
			if !noInstallDeps {
				if yes {
					installDeps = true
				} else {
					var err error
					installDeps, err = prompter.Confirm("Install dependencies now?", true)
					if err != nil {
						return err
					}
				}
			}

			opts := InitOptions{
				Name:        name,
				Module:      module,
				InstallDeps: installDeps,
				depsRunner:  runInDir,
			}

			return runInit(opts, ".")
		},
	}

	c.Flags().StringVar(&module, "module", "", "Go module path (default: github.com/example/<slug>)")
	c.Flags().BoolVarP(&yes, "yes", "y", false, "Accept all defaults without prompting")
	c.Flags().BoolVar(&noInstallDeps, "no-install-deps", false, "Skip automatic dependency installation")
	return c
}

// InitOptions holds all resolved parameters for project initialization.
type InitOptions struct {
	Name        string
	Module      string
	InstallDeps bool
	// depsRunner runs a command in a given directory; defaults to runInDir.
	depsRunner func(dir, name string, args ...string) error
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

// runInit scaffolds a new workshop project and optionally installs dependencies.
func runInit(opts InitOptions, baseDir string) error {
	slug := slugify(opts.Name)
	dest := filepath.Join(baseDir, slug)

	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("directory %q already exists", dest)
	}

	data := projectData{
		Name:   opts.Name,
		Slug:   slug,
		Module: opts.Module,
	}

	type entry struct {
		tmpl string
		dest string
	}

	files := []entry{
		{"files/project/main.go.tmpl", "main.go"},
		{"files/project/go.mod.tmpl", "go.mod"},
		{"files/project/Makefile.tmpl", "Makefile"},
		{"files/project/.gitignore.tmpl", ".gitignore"},
		{"files/project/README.md.tmpl", "README.md"},
		{"files/project/.github/workflows/release.yml.tmpl", ".github/workflows/release.yml"},
		{"files/project/site/index.html.tmpl", "site/index.html"},
		{"files/project/site/workshop.css.tmpl", "site/workshop.css"},
		{"files/project/site/workshop-progress.js.tmpl", "site/workshop-progress.js"},
		{"files/project/site/00-setup/index.md.tmpl", "site/00-setup/index.md"},
	}

	tmplFS := templates.FS()
	for _, f := range files {
		if err := scaffold(tmplFS, f.tmpl, filepath.Join(dest, f.dest), data); err != nil {
			return fmt.Errorf("scaffold %s: %w", f.dest, err)
		}
	}

	fmt.Printf("\nWorkshop %q created in %s/\n\n", opts.Name, dest)

	if opts.InstallDeps {
		runner := opts.depsRunner
		if runner == nil {
			runner = runInDir
		}
		if err := installDeps(runner, dest); err != nil {
			return err
		}
		fmt.Printf("\nRun `make serve` to start your workshop.\n")
	} else {
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", slug)
		fmt.Println("  go mod tidy")
		fmt.Println("  make serve")
	}

	return nil
}

// installDeps runs go mod tidy in the project directory.
func installDeps(run func(dir, name string, args ...string) error, dir string) error {
	fmt.Printf("  → go mod tidy\n")
	if err := run(dir, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	return nil
}
