---
name: code-review
description: >-
  Reviews local, branch, or pull-request diffs for correctness and semantics.
  Use when the user asks for a code review or writes "code-review" exactly.
---

# Code review

Review the requested diff as a senior engineer: catch real bugs, judge semantics
and idiomatic code, and give constructive feedback. Prioritize correctness over
style. Do not nitpick style for its own sake.

Coding and style standards come from **steering**, not this skill. 

## Resolve the diff

Accept **local changes** (uncommitted, staged, or branch vs base), a **ticket
key**, or a **PR URL**. If nothing is passed, review the current branch against
the updated base.

- **PR URL**: use `gh` when available, otherwise the GitHub API, to inspect the
  PR and check out its head branch (confirm first if that would move HEAD).
- **Baseline:** fetch and update `main`/`master`, or the PR destination branch.
  Compare against that ref (`origin/<base>...HEAD`). Do not use a stale local
  base.

Read every changed file and enough surrounding context. Do not review unrelated
files unless required to understand the change.

## Review workflow

### 1. Understand the change

Summarize what the change does and why (1–2 sentences). Note affected packages,
APIs, runtime paths, and assumptions the author appears to be making.

### 2. Read for correctness

- Logic errors, off-by-one, races, nil dereferences
- Incorrect error propagation or swallowed errors
- Missing edge cases (empty input, timeouts, concurrent access, partial failure)
- Breaking API or behavioral changes without migration or call-site updates
- Security-sensitive paths: auth, input validation, secrets, injection

### 3. Assess idiomatic code

Read applicable steering and judge naming, layout, interfaces, and tests against it 
and against neighboring files in the same package and layer. 
Flag pattern drift only when the change diverges from those sources. 
Do not invent team standards here, and do not suggest rewrites that fight project style.

### 4. Evaluate error handling

Treat this as a correctness pass, not a style recitation. Flag ignored or
swallowed errors, log-and-continue on paths that should fail or retry, leaked
internals in user-facing errors, and missing cleanup on error paths. What
counts as correct wrapping, context propagation, and logging comes from
steering.

### 5. Assess tests in the diff — do not run them

CI runs tests. Do not run the test suite or collect coverage locally.

Judge whether new/changed tests cover the behavior and failure paths, assert
meaningful outcomes, match project test conventions from steering, and would
catch a regression.

**Mocks:** local mocks may be stale because they are regenerated at build time
and not committed. Do not report that as an issue.


## Severity

Every finding cites `file:line` (or `file` for file-level issues) and why it
matters. Do not inflate severity or pad with speculation.

| Severity     | When to use                                                                                            |
| ------------ | ------------------------------------------------------------------------------------------------------ |
| **Critical** | Bugs, security issues, data loss, broken builds, missing error handling on production failure paths    |
| **Medium**   | Maintainability, missing tests for important paths, unclear APIs, evidenced performance, pattern drift |
| **Nitpicks** | Style, naming, optional simplifications, docs — only when worth mentioning                             |

## Report format

```markdown
# Code Review

## Summary
<1–3 sentences: what changed, overall quality, merge recommendation>

## Critical
<!-- Omit if empty -->
- **[file:line]** <finding>
  - **Why:** <impact>
  - **Suggestion:** <concrete fix>

## Medium
<!-- Omit if empty -->
- **[file:line]** <finding>
  - **Why:** <impact>
  - **Suggestion:** <concrete fix>

## Nitpicks
<!-- Omit if empty -->
- **[file:line]** <finding>
  - **Suggestion:** <optional improvement>

## Verdict
**Approve** | **Approve with suggestions** | **Request changes**

<one sentence justification>
```

## Feedback style

Direct, concise and specific. No filler text.
Prefer a small code fix over abstract advice when the fix is non-obvious. 
Empty Critical/Medium with Approve is valid when the change is solid.

## Constraints

- Review only. Do not modify code, commit, or push unless the user explicitly
  asks to apply fixes.
- Do not report issues you cannot support with code evidence.
- Stay on the diff; no unrelated refactors.
