# listopt

Shared validated list-query vocabulary for catalog list endpoints (boardgame, category, mechanism, contributor).

HTTP `GET` query parameters and `QUERY` JSON bodies parse into a transport `QueryRequest`. `toQuery(allowlist)` validates public field names with `Allowlist.ParseSort` / `ParseFilter`, collects functional `Option`s, and calls `NewQuery`. The resulting `Query` passes unchanged through `GetAll` to the database. Pagination defaults to page `1` / size `10`; `WithPagination` overrides and normalizes it.

## Pagination

- `page` / `pageSize` (GET) or `pagination` (QUERY).
- Defaults: page `1`, page size `10`.
- `pageSize` is clamped to a max of `100`. Values below 1 become the defaults.
- Out-of-range numbers are clamped; non-integers return `422`.

## Sorting

- GET: `sort=field.order` (`asc` or `desc`).
- QUERY: `"sort": { "field": "name", "order": "asc" }`.
- `field` is the public JSON name. `Allowlist.ParseSort` maps it to a DB column only when its `FieldSpec` is sortable.
- Unknown, non-sortable, or malformed sort returns `422`.
- No sort: database default (usually `id`).

Which fields are sortable: [fields.md](fields.md).

## Filtering

- GET: `filter=field.op.value`, repeatable, AND-combined. Operator is required so values may contain dots (`filter=name.like.foo.bar`).
- QUERY: `"filters": [{ "field", "op", "value" }]`. Empty `op` means `like`.
- `Allowlist.ParseFilter` validates and maps each field into a SQL-ready `Filter`.
- Operators:
  - `like` — partial string match (`ILIKE`), string fields only
  - `eq` — exact equality (string or int). On int columns the value must parse as a number
  - `lt` `le` `gt` `ge` — numeric comparison; value must parse as a number
- Int columns bind as numbers. String columns keep the raw value (`name.eq.123` matches the text `"123"`).
- Unknown field, field not filterable, or illegal operator/value returns `422`.

Which fields are filterable: [fields.md](fields.md).
