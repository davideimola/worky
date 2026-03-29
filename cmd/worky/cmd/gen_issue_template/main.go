// gen_issue_template regenerates .github/ISSUE_TEMPLATE/bug_report.md from the
// body.md.tmpl template. Run via: go generate ./cmd/worky/cmd/
package main

import (
	"log"
	"os"

	"github.com/davideimola/worky/cmd/worky/cmd"
)

func main() {
	// go generate runs with CWD = cmd/worky/cmd/
	dest := "../../../.github/ISSUE_TEMPLATE/bug_report.md"
	f, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := cmd.GenerateIssueTemplate(f); err != nil {
		log.Fatal(err)
	}
}
