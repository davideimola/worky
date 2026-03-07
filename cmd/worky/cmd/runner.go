package cmd

import (
	"os"
	"os/exec"
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

// runInDir executes a command in the given directory, streaming output.
func runInDir(dir, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
