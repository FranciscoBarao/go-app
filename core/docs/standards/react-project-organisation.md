# React project organisation

A frontend service keeps its application code under `internal/`, the same boundary the Go services use. Files are grouped by feature, and each feature exposes one public entry point.

## Root

The project root holds repository files only.

```text
frontend/
  README.md
  .gitignore       # dist/, coverage/, local env files
  Makefile
  Dockerfile
  nginx.conf       # SPA fallback for the production image only
  index.html       # Vite entry, points at internal/app/main.tsx
  package.json
  biome.json
  tsconfig*.json
  .env.example
  public/          # static files copied as-is, keeps its URL
  tests/           # cross-feature tests
  dist/            # Vite build output, gitignored
  internal/
```

`dist/` is the output of `make build`. It is never committed and is not edited by hand.

`public/` holds files that keep their URL, such as `/favicon.ico`. Vite copies them without processing.

`tests/` holds tests that cross features. A feature's own tests sit next to the code they test.

## Internal

Everything under `internal/` is application code. No other project imports it.

```text
internal/
  app/            # main.tsx, routes.tsx, app shell and nav links
  pages/          # one file per route
  features/       # one folder per feature, see below
  services/       # backend clients, one folder per service
  config/         # environment variables and app settings
  components/     # UI shared by more than one feature
  hooks/          # hooks shared by more than one page
  assets/         # images, css, fonts imported by code
  store/          # app-wide state
```

Each folder is created when its first file exists. The skeleton creates `app/`, `pages/`, `features/`, `services/`, and `config/`.

### App

`app/` holds the entry point, the router, and the shell that wraps every page.

`routes.tsx` maps a URL to a page component and nothing else. It is one module. A `routes/` folder is added only when the table splits into more than one module.

### Pages

`pages/` holds one file per route. A page composes components exported from features and shared components. It does not call the network and does not hold business logic.

A hook used only by one page sits beside that page file.

### Features

`features/<name>/` holds the code for one feature. A feature is named after the thing the user sees, in the singular, such as `boardgame` or `wishlist`.

Inside a feature, code is grouped by type. The folders mirror `internal/` without `features/` itself: `components/`, `hooks/`, `services/`, `store/`, `assets/`, `config/`. A subfolder is added when the feature has a file of that type.

`index.ts` is the feature's public API. It is the only file other code imports from the feature. Anything not exported from `index.ts` is private to the feature.

```text
internal/features/boardgame/
  index.ts          # public API: list and detail components
  components/
  hooks/
```

A feature imports from `services/`, `components/`, `hooks/`, and `config/`. It does not import from another feature. Two features that need the same code move it to `components/`, `hooks/`, or `services/`.

### Services

`services/<name>/` holds the client for one backend service, named after that service. A client holds the response models, the request functions, and the endpoint paths.

A client does not import from `features/` or `pages/`. The catalog client is `services/catalog/`. A future list service is `services/lists/`.

### Components

`components/` holds UI used by more than one feature. A component used by one feature stays in that feature.

`components/layouts/` holds layout components such as a sidebar, navbar, or container. It is added only when the first layout component exists.

### Hooks

`hooks/` holds hooks used by more than one page. A hook used by one page sits beside that page. A hook used by one feature sits in that feature.

### Store

`store/` holds app-wide state: Redux, Zustand, or a React context that is shared across the app. It is added only when that state exists.

There is no top-level `contexts/` folder. A context shared across the app is app-wide state and goes in `store/`. A context used by one feature stays inside that feature.

### Assets

`assets/` holds files imported by code. Vite bundles them and gives them hashed names in `dist/`.

`public/` and `assets/` do not overlap. `public/` keeps the URL and is copied as-is. `assets/` is processed by the bundler.

### Config

`config/` reads environment variables and exposes settings as typed values. The rest of the app imports those values and does not read `import.meta.env` directly.

Variables are prefixed with `VITE_` because Vite only exposes that prefix to the browser. A variable is baked in at build time, so a change needs a new build.

Each variable is documented in `.env.example`.

## Naming

- A file that exports a component uses PascalCase, such as `BoardgameListPage.tsx`.
- Any other file uses camelCase, such as `service.ts` or `config.ts`.
- Feature folders are lowercase and singular, such as `boardgame`. Service folders are lowercase and match the service name, such as `catalog`.
- Public entry points are `index.ts`.

## Tests

Unit tests sit next to the code they test, as `<name>.test.ts` or `<name>.test.tsx`. Tests that cross features go in `tests/`.

Tests use Vitest. `make test` runs them once. `make test/cover` writes coverage to `coverage/`.