---
name: task-manager
description: Tracks development phases, tasks, and progress via the root ROADMAP.md. Use before writing code in this repo to check what to work on next, and after doing work to update task status and log progress. Trigger terms include roadmap, tasks, PRD, progress, phase, and "what next".
---

# Task Manager

## Core rules

You are a task manager. Before writing code, always read task list and check if working on one of the tasks.
After working partially on a task, update the document to show progress / changes.
After completing a task, update the file to mark it as done.
If working on a completely different task, we should still update the task list to show progress done.

## Roadmap location

The task list lives in `ROADMAP.md` at the repo root (`/Users/fbarao/Documents/Personal/go-app/ROADMAP.md`).
This is the single source of truth for phases, tasks, and progress.

## Workflow

Before writing code:
1. Read `ROADMAP.md`.
2. Identify whether the requested work matches an existing task (or phase).
3. If it matches, mark that task in progress (`[~]`) before starting.
4. If it does not match any task, add it (see "Off-roadmap work" below) before starting.

While/after working:
1. If you made partial progress, keep the task `[~]` and append a short dated progress note.
2. If you completed the task, mark it done (`[x]`) and append a short dated note.
3. If you worked on something not on the roadmap, still record it so the list reflects reality.

## Status conventions

Use checkbox markers on each task line:
- `[ ]` — todo (not started)
- `[~]` — in progress / partial
- `[x]` — done
- Append `— Deferred` (with a reason) for intentionally postponed items.

Record progress with a nested, dated bullet under the task, e.g.:

```markdown
- [~] Pagination (page / pageSize + total count envelope)
  - 2026-07-06: added page/pageSize parsing; total count query still pending
```

## Off-roadmap work

If the current work is not on the roadmap:
1. Add it under the most relevant phase, or under an `## Ad-hoc / unplanned` section if none fits.
2. Give it a status marker and a short dated note describing what changed.

## Editing rules

- Preserve the existing phase structure and ordering; only change status markers and add notes.
- Keep notes short (one line) and dated (`YYYY-MM-DD`).
- Do not delete completed or deferred tasks — they are the progress record.
- When a whole phase is complete, mark its heading with `✅`.
