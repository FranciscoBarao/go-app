# Skills

Each directory is one Cursor skill. The agent loads `SKILL.md` when the user names the skill or the request matches its description.

A skill states a procedure. Naming, commits, and reviews still follow the standards under `core/docs/standards` and steering under `core/ai/steering`.

## Index

| Skill                                 | Use when                                                             | Standard                                                                                                    |
| ------------------------------------- | -------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| [task-manager](task-manager/SKILL.md) | Checking what to work on next, or after work that should be recorded | Root `ROADMAP.md`                                                                                           |
| [code-review](code-review/SKILL.md)   | Reviewing a local diff, a branch, or a pull request                  | Steering; not a naming standard                                                                             |
| [git-branch](git-branch/SKILL.md)     | Creating a working branch                                            | [branch-nomenclature.md](../../docs/standards/branch-nomenclature.md)                                       |
| [git-pr](git-pr/SKILL.md)             | Opening or updating a pull request                                   | [conventional-commits.md](../../docs/standards/conventional-commits.md), `.github/pull_request_template.md` |

`git-pr` delegates a missing branch to `git-branch`.

## Layout

```
skill-name/
└── SKILL.md
```

`SKILL.md` has YAML frontmatter (`name`, `description`) and the steps the agent follows. The description lists trigger terms. Do not duplicate a standard inside a skill; link it.
