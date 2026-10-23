# Persistence

Service repositories use the 3 layer architecture: transport, service, and storage. This page states how the storage layer is organised. Layer responsibilities stay in that architecture section. Error mapping stays in `error-handling.md`.

The path before `store/` depends on the repository or bounded context. Everything under `store/` is persistence.

## Layout

`store/<technology>/` identifies the database, for example `postgres` or `redis`. One technology is one package. A repository with several databases has one subdirectory per technology.

`store/<technology>/<technology>.go` holds the store struct, its constructor, its shutdown method, and other bootstrap. It does not contain model-specific logic.

`store/<technology>/<model>.go` holds the store methods for one model. Those methods are the database interactions.

SQL lives in `store/<technology>/sql/`, one file per model, named after the model. Shared fragments live in the same package. The technology package is the only importer. SQL keywords are capitalized.

## Migrations

Migrations use golang-migrate. They live in `store/<technology>/migrations/`. A migration is an up file and a down file. The name starts with a version, then a description, then `up` or `down`.

```
000001_bootstrap.up.sql
000001_bootstrap.down.sql
```

A schema change is a new migration. An already applied migration is not edited.

The migration path comes from configuration. The constructor applies pending up migrations while the store starts. No migrations to apply is success. Down files are the rollback path when a rollback is run explicitly. Startup does not run them.

## Errors

A technology-specific error is mapped to a storage sentinel before it leaves the store. Driver packages, driver types, and backend error codes stay inside the persistence layer. They do not appear in an external response. The adapter may log diagnostic context.

Storage sentinels live in the `store` package, in its error file, with no internal imports. The service can import that package even when the adapter imports the domain packages. The rules for those sentinels are in `error-handling.md`.

## Tests

Unit tests live in the technology package, next to the implementation. They do not start or connect to a real or embedded database. Database interactions are replaced with mocks or fakes.

`<model>_test.go` covers the meaningful success and error paths of that model's store methods. One model is one test file.

`<technology>_test.go` is unnecessary when it would only retest bootstrap already covered elsewhere. Add it when the repository has technology bootstrap of its own that needs unit coverage.

Integration tests live in the `tests/` folder at the root of the project. They use a separate test package and import the store from outside the package. They run against a local database, or the team-approved equivalent, and cover behavior a fake cannot prove: migrations, constraints, query scanning, and technology-specific error mapping.
