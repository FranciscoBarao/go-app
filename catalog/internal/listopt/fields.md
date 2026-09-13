# Query field allowlists

List endpoints accept only fields declared on each resource’s `QueryAllowlist`. Public names are JSON names (`min_players`, not `minplayers` or `MinPlayers`).

## Patterns

- **1 — Explicit allowlist (this package).** `map[string]FieldSpec` keyed by public name: column, kind, sortable, filterable. Lookup is a map. Sort and filter are independent flags.
- **2 — Opt-in struct tags.** Still reflection; capability lives on the domain model. Cannot express “this name is not a struct field” (joins, computed sorts) without dummy fields. `*int` / `time.Time` still leak Go kinds unless tags repeat the same data an allowlist would hold.
- **4 — Per-endpoint switch.** `switch field { case "name": ... }`. Fine for one resource; four list APIs would duplicate operator rules and drift.
- **5 — Generate from OpenAPI/SQL.** Right when the contract already lives outside Go. Overkill for four tables and hand-written handlers.

**Why 1:** four resources, a stable public API, sort ≠ filter, no request-time reflection. A new column is not queryable until it is added to the allowlist. `min_age` is `KindInt` even though the Go field is `*int`. `bgg_id` can exist on the model and stay off the query surface.

## Implementation

- [`allowlist.go`](allowlist.go): `Allowlist`, `FieldSpec`, `KindString` / `KindInt`, `Allowlist.ParseSort`, `Allowlist.ParseFilter`.
- Each domain package owns `QueryAllowlist` (`boardgame/query.go`, `category/query.go`, `mechanism/query.go`, `contributor/query.go`). `listopt` does not import domain types.
- Transport: `QueryRequest.toQuery(allowlist)` collects `Option`s, calls `NewQuery`, and returns `Query` → `GetAll`.


## Allowlists

**Boardgame** — sort and filter:

- `id`, `slug`, `name`
- `year_published`, `min_players`, `max_players`, `min_play_time`, `max_play_time`, `min_age`

Sort only: `created_at`, `updated_at`

**Category, mechanism, contributor** — sort and filter:

- `id`, `slug`, `name`

Sort only: `created_at`, `updated_at`
