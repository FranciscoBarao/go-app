# AGENTS.md - GO-APP Agent Instructions

Use these instructions in GO-APP service repositories.

## Instruction Precedence

- Follow the user's request first, then the most specific instructions in the target repository, then this file.
- Keep changes focused on the requested repository and behavior.

## Safety and Authorization

- Treat review tasks as read-only unless the user authorizes implementation.
- QA work may modify test files. Do not modify code during QA work without explicit authorization.
- Do not commit, push, deploy, apply, patch, create, or change systems without explicit authorization.

## Working Style

- Be pragmatic. Communicate concisely and focus on the task.
- After three failed attempts to solve the same problem, stop and ask for guidance. Do not enter a debug loop or improvise a workaround.


## Change Workflow

- Run the target repository's `make check` prior to committing any code; it must pass in order to commit


## Go Changes

Apply the following sections when changing Go code.

### Repository Layout and Important Repositories

- Each directory is an independent Go project with its own module, `Makefile`, `Dockerfile`.
- A typical service repository contains:
  - `config/`, `dockerfile/`, `environment/`: configuration and deployment material.
  - `build/`: generated and normally ignored artifacts.
  - `internal/`: private handlers, business logic, configuration, and adapters.
  - `pkg/`: exported packages used by other services.
  - `tests/`: integration tests, normally using the `integration` build tag.

### 3 Layer Architecture 

  - Use this architecture to understand, modify, and extend service repositories.
  - The three layers are transport, service, and storage.
  - Identify a layer by its responsibilities and dependencies, not by its package name alone.
  - Treat statements about what a layer does, accepts, returns, or depends on as architectural boundaries.
  - Treat package paths, filenames, and type names as organizational conventions:
    - When modifying an existing repository, follow its equivalent convention and avoid introducing a second structure for the same purpose.
    - When creating a new component without an established convention, use the paths and names described below.
  - Follow repository-specific guidance when a nested `AGENTS.md` documents a different architecture.
  - Do not refactor unrelated code to make it conform to this guidance. When existing code crosses a layer boundary, do not extend that coupling. Refactor it when the requested change requires the boundary to change.

#### Transport Layer Responsibilities and Conventions

  - Implement transport code under internal/transport.
  - Receive HTTP requests, perform request-shape validation, unmarshall HTTP fields into domain types, invoke the service layer, and convert results into responses.
  - Keep HTTP references within the transport layer. Do not pass them to the service or storage layer.
  - Keep business validation and authorization decisions in the service layer.
  - Define narrow interfaces for the service methods required by each transport server. Accept implementations of those interfaces through the server constructor.
  - Construct transport servers in the composition root.
  - Name client-facing servers Public<Service>Server and internal servers Private<Service>Server.
  - Translate service-layer errors into an appropriate HTTP response via middleware.ErrorHandler. 
  - Group HTTP implementations for the same domain object in <domain_object>.go.

#### Service Layer Responsibilities and Conventions

  - Implement business workflows, domain validation, authorization decisions, and coordination with storage or outbound service adapters.
  - Place each domain service under internal/<domain_object>.
  - Name the transport-facing type Service and implement it in internal/<domain_object>/service.go.
  - Define domain models in the same package as the Service methods that use them. Use model.go for a small model set or <domain_model>.go for separate models.
  - Define narrow interfaces for required storage operations and outbound services. Accept implementations through the Service constructor.
  - Construct service objects in the composition root and pass them to transport server constructors.
  - Convert storage errors into service errors that do not expose backend details. Preserve unexpected errors with %w. Storage errors must not cross the service boundary.
  - Do not log an error and return it unchanged. Log failures only when the service handles, suppresses, retries, or reports them.

#### Storage Layer Responsibilities and Conventions

  - Read and write application state through storage systems such aS PostgreSQL.
  - Own queries, serialization, storage-model mapping, backend constraints, and transactions. Do not implement business workflows, authorization policy, or transport behavior.
  - Implement new storage adapters under internal/storage/<backend>. Follow the repository’s established package convention when modifying existing code.
  - Satisfy the narrow storage interfaces defined by the consuming service package.
  - Accept and return domain models, declared result types, identifiers, or counts. Do not expose driver types, database rows, commands, or backend error codes.
  - Translate expected backend conditions into stable errors such as ErrNotFound, ErrAlreadyExists, and ErrConflict. Treat these errors as part of the storage contract and support inspection with errors.Is.
  - Wrap unexpected errors with operation and entity context using %w. Do not include credentials, tokens, connection strings, or stored payloads.
  - Enforce storage constraints such as uniqueness, serialization validity, and optimistic concurrency. Leave business validation to the service layer.
  - When writes to one backend require atomicity, expose one storage method that performs them in one transaction. Let the service coordinate operations across backends.
  - Accept context.Context for every I/O operation and pass it to the storage client.
  - Create storage clients in the composition root and close them during shutdown.
  - Name methods after entity persistence operations. Do not expose commands such as HSet, QueryRow, or Exec through service-facing interfaces.

### Formatting and Naming

- Keep lines within 160 characters where possible. When wrapping:
  1. Break after the initial package or method invocation.
  2. Break chained method invocations.
  3. Put one input argument on each line.
  4. Break at the closing bracket of a struct argument initialization.
  5. For function declarations:
    - place each argument on its own line
    - place each return type on its own line
- Run `gofmt`; rely on it to sort imports.
- Do not add an import alias that matches the imported package name.
- Keep import aliases lowercase.
- Give variable names at least one noun where possible. Receiver names, `err`, and variables scoped to a small `if`, `else`, or loop body are exceptions.
- Do not provide return variable names in function declarations
- go.mod should contain 2 `require` blocks: one for direct dependencies and a second for indirect dependencies. Direct dependencies should be listed first. 

### Errors, Context, and Logging
<!---
TODO - define once errors are improved

- Use `errors.New` instead of `fmt.Errorf` when the message contains no interpolation.
- Use formatted or wrapped errors when including values or an underlying error using `fmt` wrapping.
- Preserve context propagation, cleanup, and explicit error handling.
- Do not swallow errors or log and continue on paths that should fail.
- Prefer explicit scalar logging fields. Avoid `Interface` unless necessary.
- Prefer not combining function calls and error checking into a single if statement.
-->

### Retries and Backoff

- Before adding retry behavior, inspect the repository for an existing retry helper or client policy.
- Service repositories commonly use `github.com/cenkalti/backoff`. Use the major version already present in the repository unless the task includes a dependency migration.
- Apply retries at the outbound client or adapter boundary.
- Bound every retry policy by attempts, elapsed time, or the caller’s context.
- Stop retrying when the context is canceled.
- Retry transient failures only. Mark validation, authentication, authorization, and other terminal failures with `backoff.Permanent`.
- Confirm that an operation is idempotent before retrying it. Do not retry a mutation unless the API provides idempotency or the operation can tolerate duplicate execution.
- Keep retry timing configurable when it affects request latency or service load.
- Do not create a new repository-local retry helper when an equivalent helper already exists.

### Component Design

#### Interfaces and Constructors

- Define interfaces where consumers use them, such as the package that injects them into constructors.
- A factory may return an interface when doing so provides value.
- Use explicit positional constructor parameters for mandatory arguments and prefer functional options for optional ones.
- Inject dependencies through interfaces defined in the consuming package instead of directly naming types from another package.

#### Concurrency Ownership

- Concurrency should be the responsibility of the function caller, not the function (e.g. `go someFunc(...)` instead of `someFunc(...) { go ... }`)

#### Service Configuration 

<!-- TODO -->

### Metrics

<!-- TODO -->

### Mocks

- Never edit or create generated mocks directly. Regenerate them with `go generate`.
- Mocks are regenerated during the build and are not committed.
- If mocks are duplicated or invalid, clean and regenerate all mocks before further debugging.
- Mock code is placed in files whose name is `<filename_of_code_being_mocked_minus_dot_go>_mock.go`

### Unit Tests

- Choose table tests when testing smaller, isolated functions that have little to no setup, default to suites otherwise. Dont mix suites and table tests.
- Use `assert` for actual test assertions. Use `require` for intermediate assertions that must pass before the test can continue.
- When a test needs a generic error, use `assert.AnError`.
- When a test expects a specific error, assert that specific error.
- Prefer Gomock mocks over stubs or fakes.
- Avoid `gomock.Any()` where possible.
- Prefer exact values or `gomock.Eq()` so tests verify requests and arguments.
- Keep unit tests isolated and deterministic:
  - Do not access real networks or databases.
  - Control time and randomness.
  - Avoid sleep-based timing.
  - Clean up files, servers, ports, globals, and goroutines created by tests.
- Run lint checks with the tests.

### Test Suites

- Use `github.com/stretchr/testify/suite` when creating a unit-test suite.
- Dont include unncessary words, like `Test`, in the Suite struct name
- Name suite receivers `suite`, not `s`.
- Use `suite.Assert()` and `suite.Require()` instead of global assertion functions.
- Keep `Assert()` calls explicit; do not hide them behind helpers.
- Put the `TestXxxSuite(t *testing.T)` runner at the end of the test file.
- Do not store a `*gomock.Controller` on the suite struct.
- Prefix mock variable names with `mock`, such as `mockRepo`.

## Documentation & Responses

Apply these rules when writing documentation or comments:

- Use a professional, direct tone without becoming overly formal.
- Don´t be verbose, remove filler, pleasantries, and hedging.
- Use short words where they preserve the exact meaning.
- Prefer the pattern: `[thing] [action] [reason]. [next step].`
- Use third person, such as "the service" or "it," instead of second person.
- Do not document a method or type solely in terms of the type and package itself.
- Don´t reword the same statements in diferent ways. Keep it simple.
