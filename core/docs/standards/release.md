# Release and versioning

The repository versions with `vMAJOR.MINOR.PATCH`. The git tag is the version. There is no version file and no extra commit on `main`.

A squash-merge to `main` tags the next version and creates a GitHub Release, except when the title marks a breaking change. Major never increments from a merge. A person runs `release.yaml` with `bump=major` for that.

## Bump rules

| Source                                                                               | Result            |
| ------------------------------------------------------------------------------------ | ----------------- |
| `feat`                                                                               | minor             |
| `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert` | patch             |
| `!` after the type, or a `BREAKING CHANGE:` / `BREAKING-CHANGE:` footer              | skip auto-release |
| `workflow_dispatch` `bump=major`                                                     | major             |
| `workflow_dispatch` `bump=minor` or `bump=patch`                                     | that bump         |

The first tag is `v0.1.0` for a minor bump, `v0.0.1` for a patch bump, and `v1.0.0` for a major bump.

A skipped auto-release does not create a tag. The merge is already on `main`. The major is cut with `workflow_dispatch`.

## Workflows

The files live in `.github/workflows/` at the repository root. GitHub loads workflows only from that directory. Each file declares `workflow_call` so a later service repository can call it.

| File               | When                                                       | What it does                                                                                                                       |
| ------------------ | ---------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| `validate-pr.yaml` | pull request opened, reopened, edited, or synchronized     | Checks the PR title against conventional commits                                                                                   |
| `check.yaml`       | pull request that touches `catalog/**`, or `workflow_call` | Checks out the repo, sets up Go from `go.mod`, prints `go version`, runs `go clean -cache -modcache -testcache`, then `make check` |
| `tag.yaml`         | `workflow_call` or `workflow_dispatch`                     | Creates annotated tag `vX.Y.Z` and pushes it. A tag that already exists is left as is                                              |
| `release.yaml`     | push to `main`, `workflow_dispatch`, or `workflow_call`    | Reads the squash commit, chooses the bump, calls `tag.yaml`, creates a GitHub Release with generated notes                         |

`check.yaml` takes `working-directory`. The monorepo default is `catalog`. A caller after extract passes `.`.

The GitHub Release body is auto-generated notes. The template in [core/docs/templates/release-notes.md](../templates/release-notes.md) is for a handwritten write-up, typically a manual major. It is not the default release body.

## Manual major

Actions → Release → Run workflow → `bump=major`. That job tags the next major version and creates the GitHub Release. The release-notes template covers a change that needs a written explanation.
