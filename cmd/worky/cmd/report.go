package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"github.com/davideimola/worky/cmd/worky/templates"
	"github.com/spf13/cobra"
)

//go:generate go run ./gen_issue_template/

const reportBaseURL = "https://github.com/davideimola/worky/issues/new"

const issueTemplateFrontmatter = `---
name: Bug report
about: Report a bug in worky
title: ''
labels: bug
assignees: ''
---

`

const issueTemplateFooter = "\n> **Tip:** Run `worky report` to open this form pre-filled with your environment details.\n"

// GenerateIssueTemplate renders the GitHub issue template using body.md.tmpl
// with placeholder values. It is the single source of truth for the template structure.
func GenerateIssueTemplate(w io.Writer) error {
	placeholders := ReportInfo{
		Problem:  "<!-- Describe the problem clearly and concisely. -->",
		Steps:    "<!-- List the exact steps to reproduce the issue. -->\n\n1.\n2.\n3.",
		Expected: "<!-- What did you expect to happen? -->",
		Actual:   "<!-- What actually happened? -->",
		Version:  "<!-- e.g. v0.1.0 -->",
		OS:       "<!-- e.g. linux, darwin -->",
		Arch:     "<!-- e.g. amd64, arm64 -->",
		GoVer:    "<!-- e.g. go1.24.0 -->",
	}
	if _, err := io.WriteString(w, issueTemplateFrontmatter); err != nil {
		return err
	}
	if _, err := io.WriteString(w, reportBody(placeholders)); err != nil {
		return err
	}
	_, err := io.WriteString(w, issueTemplateFooter)
	return err
}

// NewReportCmd returns the `worky report` command.
func NewReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Report a bug in worky",
		RunE: func(cmd *cobra.Command, _ []string) error {
			env := collectEnvInfo()
			opts := ReportOptions{
				info:          env,
				browserOpener: platformBrowserOpener,
				out:           cmd.OutOrStdout(),
			}
			p := newStdinPrompter(cmd.InOrStdin(), cmd.OutOrStdout())
			return runReport(opts, p)
		},
	}
}

// platformBrowserOpener opens url using platform-specific commands.
func platformBrowserOpener(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "rundll32"
		return exec.Command(cmd, "url.dll,FileProtocolHandler", url).Start()
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}

// ReportInfo holds all information needed to build a bug report.
type ReportInfo struct {
	Problem  string
	Steps    string
	Expected string
	Actual   string
	Version  string
	OS       string
	Arch     string
	GoVer    string
}

// buildReportURL constructs a GitHub new-issue URL pre-filled with the report body.
func buildReportURL(info ReportInfo) string {
	q := url.Values{}
	q.Set("body", reportBody(info))
	return reportBaseURL + "?" + q.Encode()
}

// ReportOptions holds resolved parameters for running the report command.
type ReportOptions struct {
	info          ReportInfo
	browserOpener func(url string) error
	out           io.Writer
}

// runReport prompts the user, builds the URL, and opens the browser (or prints as fallback).
func runReport(opts ReportOptions, p Prompter) error {
	var err error
	if opts.info.Problem, err = p.Ask("Describe the problem", ""); err != nil {
		return err
	}
	if opts.info.Steps, err = p.Ask("Steps to reproduce", ""); err != nil {
		return err
	}
	if opts.info.Expected, err = p.Ask("Expected behavior", ""); err != nil {
		return err
	}
	if opts.info.Actual, err = p.Ask("Actual behavior", ""); err != nil {
		return err
	}

	u := buildReportURL(opts.info)
	if err := opts.browserOpener(u); err != nil {
		_, _ = fmt.Fprintf(opts.out, "Open this URL to submit your report:\n%s\n", u)
	}
	return nil
}

// collectEnvInfo gathers environment details for the bug report.
func collectEnvInfo() ReportInfo {
	return ReportInfo{
		Version: Version,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		GoVer:   runtime.Version(),
	}
}

func reportBody(info ReportInfo) string {
	tmplContent, err := fs.ReadFile(templates.FS(), "files/report/body.md.tmpl")
	if err != nil {
		panic("report body template missing: " + err.Error())
	}
	tmpl, err := template.New("body").Delims("<%", "%>").Parse(string(tmplContent))
	if err != nil {
		panic("report body template invalid: " + err.Error())
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, info); err != nil {
		panic("report body template execute: " + err.Error())
	}
	return buf.String()
}
