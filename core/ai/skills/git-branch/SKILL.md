---
name: git-branch
description: >-
  Creates a git branch named type/short-slug from an inferred conventional-commit
  type and slug. Use when the user asks to create a branch, start a new branch,
  git-branch, or check out a branch for the current work.
---

# Git branch

Create a working branch in the form `type/short-slug`. Format and allowed types for names are in `core/docs/standards/branch-nomenclature.md`. Scope is not used.

## Workflow

1. Infer `type` from the prompt first, then from current session changes. Use `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `ci`, or `perf`.
3. Infer the slug from the prompt first, then from current session changes. Lowercase words, hyphenate, and strip characters that are invalid in a git ref.
4. Keep the full branch name at or under 72 characters
5. Identify the base branch from repository context when possible; default to `master` when ambiguous and ask if unsure.
6. Propose the exact branch name and base. Wait for explicit approval before creating it.
7. After approval:
   - If the working tree has uncommitted changes that would be carried onto the new branch, ask before switching.
   - If a branch with that name already exists locally or on the remote, warn the user and stop. Do not overwrite or reset it.
   - Create and check out the branch from the chosen base.

## Name

`type` comes from the inference list above. `feat/` is used, not `feature/`. The branch name has no `!`.

`short-slug` is a few lowercase hyphenated words. Strip spaces, `~`, `^`, `:`, `?`, `*`, `[`, `\`, `..`, `@{`, leading or trailing `.` or `/`, and consecutive hyphens. Truncate the full `type/short-slug` name to 72 characters without ending on a hyphen.

## Guardrails

- Never create or switch branches before approval.
- Never rewrite git history unless the user explicitly asks.
- Do not create a branch from unrelated uncommitted work without calling that out first.

