package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Runner abstracts command execution for testability.
type Runner interface {
	Run(name string, args ...string) error
}

type osRunner struct{}

func (osRunner) Run(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// NewBuildCmd returns the `worky build` command.
func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build the workshop site and binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildWorkshop(".", osRunner{})
		},
	}
}

// buildWorkshop verifies the project, runs hugo, then go build.
func buildWorkshop(dir string, runner Runner) error {
	if _, err := os.Stat(filepath.Join(dir, "hugo.toml")); os.IsNotExist(err) {
		return fmt.Errorf("hugo.toml not found — are you in a worky project directory?")
	}

	if err := runner.Run("hugo"); err != nil {
		return fmt.Errorf("hugo build failed: %w", err)
	}

	binName := filepath.Base(filepath.Clean(dir))
	binPath := filepath.Join("bin", binName)

	if err := runner.Run("go", "build", "-o", binPath, "."); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	return nil
}
