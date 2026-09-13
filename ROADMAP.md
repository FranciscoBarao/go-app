# Roadmap

Development roadmap and task list for go-app. This is the single source of truth for what to work on next and progress made.

## How to use

- Task states: `[ ]` todo, `[~]` in progress / partial, `[x]` done, `— Deferred` for postponed.
- Record progress as a nested, dated bullet under a task (e.g. `- 2026-07-06: ...`).
- Before writing code, read this file and mark the relevant task in progress.
- After partial work, update the task with a dated progress note. After completing, mark it `[x]`.
- Off-roadmap work still gets logged (add it under the right phase or `## Ad-hoc / unplanned`).
- See the `task-manager` skill (`.cursor/skills/task-manager/SKILL.md`) for full conventions.

## Phases at a glance

1. Rich model & CRUD ✅
2. Browse & query
3. Catalog images
4. Frontend
5. BGG seed import
6. Accounts & Lists

Phases are listed in intended work order. Phase 3 (images) can run in parallel with the Phase 4 skeleton, and Phase 5 (import) can overlap later frontend sub-phases:

`2  ->  4.1 skeleton (+ 3 images in parallel)  ->  4.2  ->  5 import  ->  4.3-4.5  ->  6`

High-priority engineering backlog items (CI, test coverage, logging, caching) are slotted in opportunistically alongside the phases above.

## Phase 1 — Rich model & CRUD ✅

Goal: A solid domain model and slug-based API.

- [x] Slugs on boardgames, contributors, categories, mechanisms
- [x] Rich boardgame fields (players, play time, age, description, bgg_id)
- [x] Contributors + role-based credits
- [x] Categories/mechanisms: id + slug + name
- [x] Soft/hard delete (`?hard=true`, `deleted_at`)
- [x] Tags removed
- [x] No images — Deferred by design (now scheduled: see Phase 3)
- [x] Single migration (0001)
- [x] Create flow: Request → service → DTO → storage → Boardgame
- [x] Reference-only associations (taxonomy must exist first)
- [x] Slug derived from name (immutable on update)
- [x] boardgame_names / alternate names — Removed (search via ILIKE later)

### Still open polish

- [x] README / Swagger out of date
  - 2026-07-06: rewrote catalog README (slug-based API, contributors, domain model + ER diagram); regenerated Swagger (`docs.go`/`swagger.json`/`swagger.yaml`) with swag v1.8.2.
- [x] Slug collision handling
  - 2026-07-06: implemented as `409 Conflict` on duplicate slug + `400` on empty slug across boardgame/category/mechanism/contributor. No auto-suffix (rejected `catan-2` style; disambiguate via name).
- [x] Non-nullable scalar fields
  - 2026-07-06: `description`, `year_published`, `min_play_time`, `max_play_time` are now NOT NULL with zero/empty "unset" sentinel + CHECK constraints. `min_age`/`bgg_id`/`boardgame_id` stay nullable.
- [x] PATCH clearing nullable fields with JSON `null` — Resolved by design
  - 2026-07-06: made the collapsible scalars non-nullable, so explicit null-clear is unnecessary. Send `""`/`0` to unset; omit to leave unchanged.
- [ ] Soft-deleted slug reuse (partial unique index `WHERE deleted_at IS NULL`)

## Phase 2 — Browse & query (next)

Goal: Make the catalog usable for discovery without a frontend.

- [x] Pagination (page / pageSize + total count envelope)
  - 2026-07-16: implemented via HTTP `QUERY` endpoints for all four list resources (boardgame, category, mechanism, contributor). JSON body carries `pagination`/`sort`/`filters`/`include_deleted`; response is a `{data, page, pageSize, totalItems, totalPages}` envelope. Pagination normalized in `listopt` (defaults page 1 / size 10, max 100). GET list endpoints now return a plain first-page array.
- [x] `GET` list endpoints: query-param pagination/sort/filter + envelope
  - 2026-07-28: GET lists were silently truncating (first 10, bare array, no total, no way to reach page 2). They now accept `?page/pageSize/sort=field.order/filter=field.op.value` (repeatable) and return the same envelope as `QUERY`, which becomes the escape hatch for filters too complex for a URL. Both methods share one validation path (`newQueryRequestFromURL` → `toQuery`) and each controller's `list` helper. Swagger documents the GET params and envelope via per-resource `*Page` mirrors, since the pinned swag (v1.8.x) renders generics as untyped objects.
- [~] Multi-filter (category, mechanism, contributor, player count, etc.)
  - 2026-07-16: scalar multi-filter shipped via `QUERY` `filters[]` (like/eq/lt/le/gt/ge, AND-combined) on model `db` columns. Association-based filters (by category/mechanism/contributor slug) still pending.
  - 2026-09-12: scalar filters/sorts use per-resource `listopt.Allowlist` values (JSON names). Association-based filters still pending.
  - 2026-09-13: list handlers build once with `NewQuery` and pass `listopt.Query` into `GetAll` (no second merge in the service).
  - 2026-09-13: `eq` on int columns requires a parseable number; string `eq` of `"123"` stays text.
  - 2026-09-13: per-resource `QueryAllowlist` definitions locked in `query_test.go`; operator tests stay in `listopt/allowlist_test.go`.
  - 2026-09-13: clarified list-query nomenclature (`Allowlist`/`Query`, `Parse*`, `toQuery`).
- [ ] Text search: `q=` with ILIKE on name + slug
- [ ] Taxonomy browse endpoints, e.g.:
  - [ ] `GET /api/category/{slug}/boardgames`
  - [ ] `GET /api/mechanism/{slug}/boardgames`
  - [ ] `GET /api/contributors/{slug}/boardgames`
- [~] `include_deleted` for admin use (partially exists on list)
  - 2026-07-16: supported on boardgame `GET` (`?include_deleted=true`) and `QUERY` (`include_deleted` body field).

## Phase 3 — Catalog images

Goal: image support for catalog resources. Independent of Phase 2; can run in parallel with the Phase 4 skeleton.

- [ ] Decide storage strategy (external URL vs object storage e.g. MinIO/S3) — decide first, it shapes the rest
- [ ] Add image fields + migration (cover + thumbnail on boardgame; optional on contributor/category)
- [ ] API accepts/returns image references; validation
- [ ] Multiple images / gallery — later
- [ ] Frontend rendering — handled in Phase 4 (4.2)

## Phase 4 — Frontend

Goal: public browse UI on top of Phase 2 APIs. Built in `frontend/` (currently a stub). No user accounts in scope here — that moves to Phase 6.

### Phase 4.1 — Skeleton (HIGH — before Phase 5)

- [ ] Framework choice + project scaffold
- [ ] Routing + base layout / navigation
- [ ] API client + environment/config wiring
- [ ] Dockerized dev setup (compose integration)

### Phase 4.2 — Browse & detail

- [ ] Boardgame list page
- [ ] Game detail page (by slug)
- [ ] Image rendering (cover/thumbnail from Phase 3)

### Phase 4.3 — Search & filter UI

- [ ] Text search (`q=`)
- [ ] Multi-filter controls (category, mechanism, contributor, player count)
- [ ] Pagination controls

### Phase 4.4 — Taxonomy pages

- [ ] Category browse page
- [ ] Mechanism browse page
- [ ] Contributor browse page

### Phase 4.5 — Polish

- [ ] Loading / error / empty states
- [ ] Responsive layout
- [ ] Basic accessibility pass

## Phase 5 — BGG seed import

Goal: Bootstrap real data. Deprioritized below the Phase 4.1 skeleton — schedule after the frontend skeleton; may run in parallel with later frontend sub-phases.

- [ ] `catalog/cmd/bggimport` CLI
- [ ] Curated seed list + `make seed-bgg`
- [ ] Map BGG XML → catalog model (bgg_id, contributors, taxonomy)
- [ ] Alternate names can wait until search needs them (or denormalized `search_text`)

## Phase 6 — Accounts & Lists

Goal: introduce user accounts and personal/group game lists. Build on the existing `user-management/` and `list-service/` scaffolding.

### Phase 6.1 — Auth foundation

- [ ] Finish user-management: register / login / JWT issuance
- [ ] Shared auth middleware across services
- [ ] Protect catalog write endpoints (currently unprotected)
- [ ] User model + roles (admin vs user)

### Phase 6.2 — Lists core

- [ ] Owned games + wishlist
- [ ] `list-service` CRUD referencing catalog by slug
- [ ] "Discontinued" handling per `list-service/README.md` decision (append-only, flag on catalog lookup)

### Phase 6.3 — Play sessions

- [ ] Log plays (game, date, players, scores)

### Phase 6.4 — Groups

- [ ] Group membership
- [ ] Aggregate "available games" across a group's owned collections

### Cross-cutting

- [ ] Service-to-service auth / API gateway considerations

## Future work — prioritized backlog

Prioritization uses priority buckets (High / Medium / Low / Very Low) mapped to MoSCoW (H = Must/Should, M = Should/Could, L = Could, VL = Won't-for-now). SWOT is intentionally not used — it is a strategic-analysis tool, not a backlog prioritizer.

| Item | Type | Priority | MoSCoW | Notes |
|------|------|----------|--------|-------|
| Caching | Eng | High | Must/Should | Response/query cache (e.g. Redis) + invalidation strategy |
| Test coverage | Eng | High | Must/Should | Evaluate how coverage is measured, then improve + add gates |
| Logging improvement / docs | Eng | High | Must/Should | Consistent structured logs, log levels, documentation |
| Development pipelines / CI | Eng | High | Must/Should | Per-service build, test, lint on PRs |
| Kubernetes test env for CD | Eng | Medium | Should/Could | Hosted test environment for continuous deployment |
| Performance / load tests | Eng | Medium | Should/Could | Throughput and latency under load |
| UI testing | Eng | Medium | Should/Could | Frontend component + e2e tests |
| Automatic releases / versioning | Eng | Low | Could | Semver, changelogs, tags |
| PR automation | Eng | Low | Could | Auto squash-merge + validate PR title & branch nomenclature |
| Game editions/versions as first-class entities | Product | Low | Could | Distinct from name aliases |
| Community (reviews, forums, ratings) | Product | Very Low | Won't (now) | Deferred |
| Shopping / marketplace | Product | Very Low | Won't (now) | Deferred |
