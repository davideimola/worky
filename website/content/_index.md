---
title: "Worky — Build Interactive Workshops as a Single Go Binary"
layout: hextra-home
---

{{< hextra/hero-badge >}}
  <div class="hx:w-2 hx:h-2 hx:rounded-full hx:bg-primary-400"></div>
  <span>Go library &rarr;</span>
{{< /hextra/hero-badge >}}

<div class="hx:mt-6 hx:mb-6">
{{< hextra/hero-headline >}}
  Build self-contained&nbsp;<br class="hx:sm:block hx:hidden" />interactive workshops
{{< /hextra/hero-headline >}}
</div>

<div class="hx:mb-12">
{{< hextra/hero-subtitle >}}
  One binary. Embedded docs. Progressive chapter locking.&nbsp;<br class="hx:sm:block hx:hidden" />Participants download one file, run one command, and follow along.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx:mb-6">
{{< hextra/hero-button text="Get Started" link="docs/getting-started" >}}
{{< hextra/hero-button text="View Demo" link="showcase/demo" style="outline" >}}
</div>

<div class="hx:mt-6"></div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Single Binary"
    subtitle="The entire workshop — docs site, checks, progress tracker — ships as one Go binary. No Docker, no cloud, no separate web server."
    icon="cube"
  >}}
  {{< hextra/feature-card
    title="Progressive Chapter Locking"
    subtitle="Chapters unlock only when the previous one passes validation. Participants can't skip ahead, keeping everyone in sync."
    icon="lock-closed"
  >}}
  {{< hextra/feature-card
    title="Live Progress Tracking"
    subtitle="The browser sidebar updates in real-time via SSE as checks pass. No refresh needed."
    icon="chart-bar"
  >}}
  {{< hextra/feature-card
    title="Embedded Site"
    subtitle="Embed HTML and Markdown files directly in your binary. The server renders Markdown on the fly — no build step, no external tools. Or skip chapter content and use the built-in progress UI when slides live elsewhere."
    icon="document-text"
  >}}
  {{< hextra/feature-card
    title="Pre-built Checks"
    subtitle="FileExists, EnvVarSet, CommandSucceeds, HTTPStatus and more — ready to use from the checks sub-package."
    icon="check-circle"
  >}}
  {{< hextra/feature-card
    title="CLI Scaffolding"
    subtitle="worky init scaffolds a complete workshop in seconds — site included. worky new chapter adds chapters with the Go snippet ready to copy."
    icon="terminal"
  >}}
{{< /hextra/feature-grid >}}
