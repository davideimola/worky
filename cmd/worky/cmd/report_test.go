package cmd

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestIssueTemplateIsUpToDate(t *testing.T) {
	var buf bytes.Buffer
	if err := GenerateIssueTemplate(&buf); err != nil {
		t.Fatalf("generateIssueTemplate: %v", err)
	}
	actual, err := os.ReadFile("../../../.github/ISSUE_TEMPLATE/bug_report.md")
	if err != nil {
		t.Fatalf("read issue template: %v", err)
	}
	if buf.String() != string(actual) {
		t.Error("issue template is out of date; run: go generate ./cmd/worky/cmd/")
	}
}

func TestRootCmd_HasReportSubcommand(t *testing.T) {
	root := NewRootCmd()
	for _, sub := range root.Commands() {
		if sub.Use == "report" {
			return
		}
	}
	t.Error("expected 'report' subcommand to be registered in root")
}

func TestBuildReportURL_IsGitHubNewIssueURL(t *testing.T) {
	info := ReportInfo{}
	url := buildReportURL(info)
	if !strings.HasPrefix(url, "https://github.com/davideimola/worky/issues/new") {
		t.Errorf("expected GitHub new-issue URL, got: %s", url)
	}
}

func TestRunReport_OpensBrowser(t *testing.T) {
	var opened string
	opts := ReportOptions{
		info: ReportInfo{Problem: "crash", Version: "v0.1.0", OS: "linux", Arch: "amd64", GoVer: "go1.24.0"},
		browserOpener: func(u string) error {
			opened = u
			return nil
		},
		out: &strings.Builder{},
	}
	if err := runReport(opts, &stubPrompter{}); err != nil {
		t.Fatalf("runReport: %v", err)
	}
	if !strings.HasPrefix(opened, "https://github.com/davideimola/worky/issues/new") {
		t.Errorf("expected browser to open GitHub URL, got: %s", opened)
	}
}

func TestRunReport_PrintsURLWhenBrowserFails(t *testing.T) {
	out := &strings.Builder{}
	opts := ReportOptions{
		info: ReportInfo{Problem: "crash", Version: "v0.1.0", OS: "linux", Arch: "amd64", GoVer: "go1.24.0"},
		browserOpener: func(u string) error {
			return fmt.Errorf("no browser available")
		},
		out: out,
	}
	if err := runReport(opts, &stubPrompter{}); err != nil {
		t.Fatalf("runReport: %v", err)
	}
	if !strings.Contains(out.String(), "https://github.com/davideimola/worky/issues/new") {
		t.Errorf("expected URL in output, got: %s", out.String())
	}
}

func TestRunReport_PromptsCollectFields(t *testing.T) {
	var openedURL string
	answers := map[string]string{
		"Describe the problem":          "app crashes",
		"Steps to reproduce":            "run serve",
		"Expected behavior":             "starts fine",
		"Actual behavior":               "panics",
	}
	p := &recordingPrompter{answers: answers}

	opts := ReportOptions{
		info:          ReportInfo{Version: "v0.1.0", OS: "linux", Arch: "amd64", GoVer: "go1.24.0"},
		browserOpener: func(u string) error { openedURL = u; return nil },
		out:           &strings.Builder{},
	}
	if err := runReport(opts, p); err != nil {
		t.Fatalf("runReport: %v", err)
	}

	parsed, _ := url.Parse(openedURL)
	body := parsed.Query().Get("body")
	for _, want := range []string{"app crashes", "run serve", "starts fine", "panics"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q\nbody:\n%s", want, body)
		}
	}
}

// stubPrompter returns empty strings for all questions.
type stubPrompter struct{}

func (s *stubPrompter) Ask(_, _ string) (string, error)       { return "", nil }
func (s *stubPrompter) Confirm(_ string, d bool) (bool, error) { return d, nil }

// recordingPrompter answers based on a map keyed by question prefix.
type recordingPrompter struct {
	answers map[string]string
}

func (r *recordingPrompter) Ask(question, _ string) (string, error) {
	for prefix, ans := range r.answers {
		if strings.Contains(question, prefix) {
			return ans, nil
		}
	}
	return "", nil
}

func (r *recordingPrompter) Confirm(_ string, d bool) (bool, error) { return d, nil }

func TestCollectEnvInfo_PopulatesOSAndArch(t *testing.T) {
	info := collectEnvInfo()
	if info.OS == "" {
		t.Error("expected non-empty OS")
	}
	if info.Arch == "" {
		t.Error("expected non-empty Arch")
	}
	if info.GoVer == "" {
		t.Error("expected non-empty GoVer")
	}
}

func TestBuildReportURL_ContainsProblemAndEnvironment(t *testing.T) {
	info := ReportInfo{
		Problem: "the server crashes on startup",
		Steps:   "run worky serve",
		Version: "v1.2.3",
		OS:      "linux",
		Arch:    "amd64",
		GoVer:   "go1.24.0",
	}
	raw := buildReportURL(info)

	// parse back the body param
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("invalid URL: %v", err)
	}
	body := parsed.Query().Get("body")

	for _, want := range []string{
		"the server crashes on startup",
		"run worky serve",
		"v1.2.3",
		"linux",
		"amd64",
		"go1.24.0",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q\nbody:\n%s", want, body)
		}
	}
}
