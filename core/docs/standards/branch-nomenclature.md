# Branch nomenclature

A working branch uses the same types as [conventional commits](conventional-commits.md). The name is a short slug for local work and review. The squash commit is the pull request title, not the branch name.

There is no ticket in the name. Scope is not used.

## Format

```
type/short-slug
```

`type` is one of `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, or `revert`. `short-slug` is lowercase words separated by hyphens. It names the change in a few words, not the full pull request title.

`feat/` is used, not `feature/`. A breaking change is marked on the pull request title (`feat!: …`). The branch name has no `!`.

One branch per pull request. The branch is deleted after the squash-merge.

Branch names are not linted in CI.

## Examples

Valid:

- `feat/boardgame-list-filters`
- `fix/empty-slug`
- `docs/conventional-commits`
- `ci/make-check`
- `refactor/drop-unused-tag-field`

Invalid:

- `feature/boardgame-list-filters` (unknown type)
- `feat/catalog/boardgame-list-filters` (scope)
- `feat/add Boardgame List Filters` (spaces, capitals)
- `feat!/drop-field` (`!` in the name)
