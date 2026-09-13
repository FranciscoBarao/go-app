# Query field schemas

List endpoints accept only fields declared on each resource’s `QuerySchema`. Public names are JSON names (`min_players`, not `minplayers` or `MinPlayers`).

## Patterns

- **1 — Explicit schema (this package).** `map[string]Field` keyed by public name: column, kind, sortable, filterable. Lookup is a map. Sort and filter are independent flags.
- **2 — Opt-in struct tags.** Still reflection; capability lives on the domain model. Cannot express “this name is not a struct field” (joins, computed sorts) without dummy fields. `*int` / `time.Time` still leak Go kinds unless tags repeat the same data a schema would hold.
- **4 — Per-endpoint switch.** `switch field { case "name": ... }`. Fine for one resource; four list APIs would duplicate operator rules and drift.
- **5 — Generate from OpenAPI/SQL.** Right when the contract already lives outside Go. Overkill for four tables and hand-written handlers.

**Why 1:** four resources, a stable public API, sort ≠ filter, no request-time reflection. A new column is not queryable until it is added to the schema. `min_age` is `KindInt` even though the Go field is `*int`. `bgg_id` can exist on the model and stay off the query surface.

## Implementation

- [`schema.go`](schema.go): `Schema`, `Field`, `KindString` / `KindInt`, `Schema.Sort`, `Schema.Filter`.
- Each domain package owns `QuerySchema` (`boardgame/query.go`, `category/query.go`, `mechanism/query.go`, `contributor/query.go`). `listopt` does not import domain types.
- Transport: `QueryRequest.toOptions(schema)` → `[]listopt.Option` → `Apply` → `Params` → `GetAll`.

## Tests

- Operator and kind rules (like vs eq, int vs string, numeric `eq`): [`schema_test.go`](schema_test.go).
- Per-resource allowlists (keys, columns, sort vs filter, omit-list): `boardgame/query_test.go`, `category/query_test.go`, `mechanism/query_test.go`, `contributor/query_test.go`.
- HTTP omitted-field smoke: `internal/transport/query_test.go` (`TestQueryRequest_ToOptions_OmittedFields`).

Omitted on every resource: `deleted_at` (use `include_deleted` on boardgames), associations, `boardgame_id`, `description`, `bio`, `bgg_id` (external BoardGameGeek id, not a browse key).

Timestamps (`created_at`, `updated_at`) are sortable only; timestamp filter binding is not implemented.

## Allowlists

**Boardgame** — sort and filter:

- `id`, `slug`, `name`
- `year_published`, `min_players`, `max_players`, `min_play_time`, `max_play_time`, `min_age`

Sort only: `created_at`, `updated_at`

**Category, mechanism, contributor** — sort and filter:

- `id`, `slug`, `name`

Sort only: `created_at`, `updated_at`
