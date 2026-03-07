package templates

import (
	"embed"
	"io/fs"
)

//go:embed all:files
var templateFS embed.FS

// FS returns the embedded template filesystem.
// Templates use [[ ]] as delimiters to avoid conflicts with Hugo shortcode syntax {{ }}.
func FS() fs.FS {
	return templateFS
}
