---
name: epic-to-issues
description: Break down a complex GitHub issue (epic) into structured, independently-mergeable sub-issues with proper labels, blocking relationships, and a tasklist on the parent. Use this skill whenever the user wants to split a large or complex GitHub issue into smaller tasks, mentions "epic", "sub-issues", "break down this issue", "spezzare questa issue", or wants to structure work so multiple contributors can work in parallel without long-lived branches. Also trigger when the user references /prd-to-issues in the context of an existing GitHub issue.
---

# Epic to Issues

Break down a large GitHub issue (epic) into structured, independently-mergeable sub-issues. The goal is trunk-based development: each sub-issue becomes a PR that goes to `main` without breaking anything, usually because new behavior is opt-in or behind a flag.

## Process

### Step 1 — Read the epic

Fetch the issue with `gh issue view <number>` and read it carefully. Identify the natural boundaries: what are the distinct components, commands, or modules mentioned? What depends on what?

### Step 2 — Propose a breakdown (collaborate, don't just execute)

Present a table of proposed sub-issues to the user before creating anything. For each item, note:
- What it delivers
- Whether it's **Core MVP** or **Beta** (useful but not blocking)
- Whether it can be merged independently to `main`

Ask the user to review: should any be merged together? Split further? Reordered? Moved between MVP and beta? Listen carefully — the user knows the codebase and the contributor audience.

Key questions to resolve before proceeding:
- Which issues can be worked in parallel? (same dependency = parallel)
- Are there any that are really "nice to have" vs truly needed for the feature to be usable?
- Should tests be separate issues or part of each issue? (default: integrated, not separate)

### Step 3 — Create the sub-issues

Once the user confirms the breakdown, create each sub-issue with `gh issue create`. Each issue body should include:

```markdown
## Parent Epic
Part of #<N> — <epic title>

## Scope
<1-2 sentences describing what this issue delivers>

## Acceptance Criteria
- <concrete, testable outcome>
- ...

## Testing
- <what to verify, following existing test patterns in the codebase>

## Dependencies
- #<N> (<title>) — if blocked by another issue
```

For **beta** issues, add a `## Status` section after the parent epic line:
```markdown
## Status
**Beta** — <MVP issues> must be completed first. In the meantime, users can <manual workaround>.
```

### Step 4 — Create labels if needed

Before creating issues, check existing labels with `gh label list`. Create any missing labels:
- `epic` (color `#6B46C1`) — for the parent issue
- `beta` (color `#FFA500`) — for future/beta sub-issues

Apply `epic` to the parent. Apply `beta` to beta sub-issues.

### Step 5 — Link sub-issues natively

Use the GitHub GraphQL API to add each sub-issue to the parent:

```bash
gh api graphql -f query='
  mutation($parentId: ID!, $childId: ID!) {
    addSubIssue(input: {issueId: $parentId, subIssueId: $childId}) {
      issue { number }
      subIssue { number }
    }
  }' \
  -f parentId="$(gh issue view <parent> --json id --jq '.id')" \
  -f childId="$(gh issue view <child> --json id --jq '.id')"
```

### Step 6 — Add blocking relationships

Use `addBlockedBy` to express dependencies. The pattern: issue A blocks issue B means B cannot start until A is done.

```bash
gh api graphql -f query='
  mutation($issueId: ID!, $blockingIssueId: ID!) {
    addBlockedBy(input: {issueId: $issueId, blockingIssueId: $blockingIssueId}) {
      clientMutationId
    }
  }' \
  -f issueId="<ID of the blocked issue>" \
  -f blockingIssueId="<ID of the blocking issue>"
```

Get IDs with: `gh issue view <number> --json id --jq '.id'`

If multiple issues are blocked by the same one (e.g., a foundational module blocks everything else), loop through them:

```bash
BLOCKING_ID=$(gh issue view <N> --json id --jq '.id')
for issue in <n1> <n2> <n3>; do
  BLOCKED_ID=$(gh issue view $issue --json id --jq '.id')
  gh api graphql -f query='...' -f issueId="$BLOCKED_ID" -f blockingIssueId="$BLOCKING_ID"
done
```

### Step 7 — Update the parent epic

Append a tracking section to the parent issue body. Read the current body first, append, then write back:

```bash
gh issue view <N> --json body --jq '.body' > /tmp/epic_body.txt
cat >> /tmp/epic_body.txt << 'EOF'

---

## Tracking

### Core MVP
- [ ] #<N> <title>
- [ ] #<N> <title>

### Beta
- [ ] #<N> <title>
- [ ] #<N> <title>
EOF
gh issue edit <N> --body "$(cat /tmp/epic_body.txt)"
```

GitHub renders these checkboxes as a progress bar on the epic.

## Principles

**Trunk-based first.** Every sub-issue should be mergeable to `main` without breaking existing behavior. Opt-in flags, new commands, and additive API changes are all safe. If a sub-issue would break something on merge, flag it and discuss.

**Tests are part of the issue, not a separate one.** Each issue's acceptance criteria should include what to test and which existing test patterns to follow. Don't create a separate "testing" issue unless the test infrastructure itself is the deliverable.

**Parallel where possible.** If two issues share the same dependency (both blocked by #X), they can be worked in parallel. Make this explicit in the blocking relationships and in your proposal.

**Beta = usable workaround exists.** An issue is "beta" when the feature is useful but users can reasonably work around it manually in the meantime (e.g., `docker stop` manually instead of a `done` command). Call out the workaround in the issue body.

**Don't create issues for things that are implicit.** If every issue will obviously need to follow the existing code style or use the existing test runner, don't add that as a criterion — it's noise.
