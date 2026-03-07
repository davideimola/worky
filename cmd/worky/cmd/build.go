package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewBuildCmd returns the `worky build` command.
func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build the workshop binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildWorkshop(".", osRunner{})
		},
	}
}

// buildWorkshop runs go build for the workshop.
func buildWorkshop(dir string, runner Runner) error {
	binName := filepath.Base(filepath.Clean(dir))
	binPath := filepath.Join("bin", binName)

	if err := runner.Run("go", "build", "-o", binPath, "."); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	return nil
}
