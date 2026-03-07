package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Prompter abstracts interactive user prompts for testability.
type Prompter interface {
	Ask(question, defaultVal string) (string, error)
	Confirm(question string, defaultYes bool) (bool, error)
}

// stdinPrompter reads answers from an io.Reader (typically os.Stdin).
type stdinPrompter struct {
	scanner *bufio.Scanner
	out     io.Writer
}

func newStdinPrompter(in io.Reader, out io.Writer) *stdinPrompter {
	return &stdinPrompter{
		scanner: bufio.NewScanner(in),
		out:     out,
	}
}

func (p *stdinPrompter) Ask(question, defaultVal string) (string, error) {
	if defaultVal != "" {
		_, _ = fmt.Fprintf(p.out, "? %s [%s]: ", question, defaultVal)
	} else {
		_, _ = fmt.Fprintf(p.out, "? %s: ", question)
	}
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return "", err
		}
		return defaultVal, nil
	}
	answer := strings.TrimSpace(p.scanner.Text())
	if answer == "" {
		return defaultVal, nil
	}
	return answer, nil
}

func (p *stdinPrompter) Confirm(question string, defaultYes bool) (bool, error) {
	hint := "[Y/n]"
	if !defaultYes {
		hint = "[y/N]"
	}
	_, _ = fmt.Fprintf(p.out, "? %s %s ", question, hint)
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return false, err
		}
		return defaultYes, nil
	}
	answer := strings.TrimSpace(strings.ToLower(p.scanner.Text()))
	if answer == "" {
		return defaultYes, nil
	}
	return answer == "y" || answer == "yes", nil
}

// yesPrompter accepts all prompts with their default value (used with --yes flag).
type yesPrompter struct {
	out io.Writer
}

func (p *yesPrompter) Ask(question, defaultVal string) (string, error) {
	if defaultVal != "" {
		_, _ = fmt.Fprintf(p.out, "? %s [%s]: %s\n", question, defaultVal, defaultVal)
	} else {
		_, _ = fmt.Fprintf(p.out, "? %s: \n", question)
	}
	return defaultVal, nil
}

func (p *yesPrompter) Confirm(question string, _ bool) (bool, error) {
	_, _ = fmt.Fprintf(p.out, "? %s [Y/n]: y\n", question)
	return true, nil
}
