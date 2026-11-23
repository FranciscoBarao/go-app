# Frontend

React + TypeScript app. 


## Requirements

- Node 26 (current LTS). Vite 8.3 and React 19.3 are the versions in `package.json`.
- Docker, for `make build-docker` and `make run`

## Commands

| Command             | Action                                                               |
| ------------------- | -------------------------------------------------------------------- |
| `make dev`          | Starts the Vite dev server on the host                               |
| `make lint`         | Runs Biome                                                           |
| `make build`        | Typechecks and builds into `dist/`                                   |
| `make test`         | Runs Vitest once                                                     |
| `make test/cover`   | Runs Vitest with coverage into `coverage/`                           |
| `make check`        | Runs `clean build test lint`                                         |
| `make build-docker` | Builds the production nginx image                                    |
| `make run`          | Rebuilds and restarts the `frontend` service in compose on port 5173 |

Run `make help` for the full list.

## Layout

The layout follows [core/docs/standards/react-project-organisation.md](../core/docs/standards/react-project-organisation.md).
