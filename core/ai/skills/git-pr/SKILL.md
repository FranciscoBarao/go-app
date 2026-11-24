---
name: git-pr
description: >-
  Opens or updates a GitHub pull request for the current branch. Use when the
  user asks to create a PR, open a pull request, git-pr, or draft PR title and
  body.
---

# Git PR

Open a pull request for the current work. Title format is in `core/docs/standards/conventional-commits.md`. Body follows `.github/pull_request_template.md`.

## Workflow

1. Gather PR context before drafting:
   - Check git status, branch name, upstream tracking, recent commits, and the diff from the base branch.
   - If branch does not exist yet, create it to track the local staged/unstaged changes. Delegate to `git-branch` skill
   - Include session context from the current conversation when relevant.
   - Identify the base branch from repository context when possible; default to `master` when ambiguous and ask if unsure.
   - Check whether the branch needs to be pushed.
   - Do not include unrelated uncommitted changes in the PR summary.

2. If a pull request already exists for the current branch, warn the user with the existing PR URL and provide the option of updating the PR description based on the latest state.

3. Draft the PR title using Conventional Commit style with the ticket as scope when known:

4. Draft the PR body following the PR template on the repo. Fallback to this template and warn the user if missing:

5. Push the branch with upstream tracking if needed.
   - Create the PR with `gh` when available, otherwise the GitHub API.
   - Return the PR URL.

## Guardrails

- Keep the title under 72 characters.
- Never push or create a PR from unrelated changes without calling them out first.
- Never rewrite git history unless the user explicitly asks.
- If there are local uncommitted changes that should be included, ask whether to commit them first unless the user already requested a commit.
- If upstream is unset, set it on the first push with `git push -u origin HEAD`.
  Do not retry a bare `git push`. For any other push failure, stop and report
  the error; do not create the PR.
- Never merge the PR

## Title

Follow `core/docs/standards/conventional-commits.md`. Use `type(ticket): subject` when a ticket is known. Omit the scope when it is not. Mark a breaking change with `!` after the type (`feat!: …`) and check Breaking change in the body. Infer `type` and subject from the prompt first, then from the branch name and the diff.

## Body

Read `.github/pull_request_template.md` or `.github/PULL_REQUEST_TEMPLATE.md`. Fill Summary, Problem, and Proposed changes. Leave HTML comments out of the published body. Check Breaking change only when the title has `!`.

If no template file exists, warn and use:

```markdown
## Summary

- [ ] Breaking change

## Problem

## Proposed changes
```

Do not mention an AI agent in the title or body.
