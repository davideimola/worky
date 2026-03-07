# Finishing Up

🎉 Almost there — this is the final chapter!

> This chapter is written in **Markdown** to demonstrate that worky chapters don't have to be HTML files. Just drop an `index.md` in the chapter directory and worky renders it automatically.

## What you'll do

Create a Markdown file called `complete.md` with a specific heading.
This simulates writing a completion log or a deployment summary — a common
pattern in real workshop scenarios.

## Checks

- **complete.md exists** — Create the file in the current directory.
- **complete.md contains `# Workshop Complete`** — The exact heading must appear in the file.

## Steps

```sh
# Create the completion file
echo "# Workshop Complete" > complete.md

# Run the final checks
go run . check
```

When the check passes, the entire workshop is complete. The banner at the bottom
of your browser will confirm it, and `go run . status` will show all ✅.

---

> **What happens when you complete the workshop?**
> Progress is stored in `~/.worky-demo/progress.json`.
> Run `go run . reset` to clear it and start over, or `go run . status`
> to see the current state at any time.

---

Now use worky to build your own. Add chapters, write checks, embed your site — and write your chapters in HTML or Markdown, whatever suits you best. 🏆
