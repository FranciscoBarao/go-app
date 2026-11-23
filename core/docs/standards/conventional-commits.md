# Conventional commits

Pull request titles follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) spec. The repository squash-merges with the PR title as the commit subject and an empty body, so the title is the commit that lands on `main`.

Branch commits are not linted. `validate-pr.yaml` checks the title on `opened`, `reopened`, `edited`, and `synchronize`.

## Format

```
type(optional-scope): subject
```

`type` is one of `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, or `revert`. Scope is optional. Subject is a short description on the same line.

A `!` after the type (or `type(scope)`) marks a breaking change. A `BREAKING CHANGE:` or `BREAKING-CHANGE:` footer marks one as well. That marker does not bump the major version. A breaking change is released with `workflow_dispatch` and `bump=major`. See [release.md](release.md).

## Examples

Valid:

- `feat: add boardgame list filters`
- `feat(catalog): add boardgame list filters`
- `fix: reject an empty slug`
- `docs: describe conventional commits`
- `refactor!: drop the unused tag field`

Invalid:

- `Add boardgame list filters` (missing type)
- `feat add filters` (missing colon)
- `feature: add filters` (unknown type)

## Enforcement

`.github/workflows/validate-pr.yaml` runs [`amannn/action-semantic-pull-request`](https://github.com/amannn/action-semantic-pull-request) on every pull request. Scope is not required.

Squash settings that make the merge commit match this title are in [release.md](release.md).
