package main

import (
	"embed"

	"github.com/davideimola/worky"
	"github.com/davideimola/worky/checks"
)

//go:embed all:site
var site embed.FS

func main() {
	worky.New(worky.Config{
		Name:    "Worky Demo",
		HomeDir: ".worky-demo",
		Port:    8080,
		SiteFS:  site,
		Chapters: []worky.Chapter{
			{
				ID:   "00",
				Name: "Getting Started",
				Slug: "00-getting-started",
				Checks: []worky.Check{
					{
						Description: "done.txt exists in current directory",
						Run:         checks.FileExists("done.txt"),
					},
					{
						Description: "WORKSHOP_USER environment variable is set",
						Run:         checks.EnvVarSet("WORKSHOP_USER"),
					},
				},
			},
			{
				ID:   "01",
				Name: "Hello Worky",
				Slug: "01-hello-worky",
				Checks: []worky.Check{
					{
						Description: "hello.txt exists",
						Run:         checks.FileExists("hello.txt"),
					},
					{
						Description: `hello.txt contains "Hello, Worky!"`,
						Run:         checks.FileContains("hello.txt", "Hello, Worky!"),
					},
				},
			},
			{
				ID:   "02",
				Name: "Finishing Up",
				Slug: "02-finishing-up",
				Checks: []worky.Check{
					{
						Description: "complete.md exists",
						Run:         checks.FileExists("complete.md"),
					},
					{
						Description: `complete.md contains "# Workshop Complete"`,
						Run:         checks.FileContains("complete.md", "# Workshop Complete"),
					},
				},
			},
		},
	}).Run()
}
