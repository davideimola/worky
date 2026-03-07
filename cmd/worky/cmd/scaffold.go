package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

// scaffold reads a template from tmplFS at tmplPath, executes it with data,
// and writes the result to destPath, creating parent directories as needed.
// Templates use [[ ]] delimiters to avoid conflicts with Hugo shortcode syntax.
func scaffold(tmplFS fs.FS, tmplPath, destPath string, data any) error {
	tmplContent, err := fs.ReadFile(tmplFS, tmplPath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Delims("<%", "%>").Parse(string(tmplContent))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
