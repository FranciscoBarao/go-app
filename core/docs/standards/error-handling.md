# Error handling

Service repositories use the 3 layer architecture: transport, service, and storage. This page states how errors cross those boundaries. Layer responsibilities stay in that architecture section.

The transport layer keeps HTTP inside itself. The handler for one resource maps that resource's errors to a response at the call site. The service does not log an error it returns unchanged. An unexpected storage failure does not cross the service boundary. The storage layer maps an expected backend condition to a stable sentinel, inspected with `errors.Is`. It does not expose a driver type or a backend error code.

## Creating an error

A new error with no cause uses `errors.New`. A wrapped error uses `fmt.Errorf("message: %w", err)`.

The message names the entry that failed, such as a slug, id, or field, when that value is in scope. It omits credentials, tokens, connection strings, SQL, and driver text. `Error()` is the client message, so the sentinel text must stay safe to show.

A driver or network error is not wrapped. Wrapping it puts driver text on the chain. The adapter logs that error, then returns `errors.New` with a safe message.

## Sentinels

| Sentinel               | Meaning                                                                                                                                                     |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ErrNotFound`          | No matching record                                                                                                                                          |
| `ErrAlreadyExists`     | The entry is already present. The name stays the same from storage through transport                                                                        |
| `ErrForeignKey`        | A foreign key failed. Storage uses this name. The service names the direction: a missing reference on write, or a record that is still referenced on delete |
| `ErrValidation`        | The request fails validation                                                                                                                                |
| `ErrMissingReference`  | A write references a record that does not exist                                                                                                             |
| `ErrStillReferenced`   | A delete is rejected because another record still points at this one                                                                                        |
| `ErrParentIsExpansion` | An expansion cannot be a parent                                                                                                                             |

A sentinel name states the condition. A foreign key failure is a missing reference on write and a record still in use on delete, so those are different sentinels.

Callers use `errors.Is`. There is no error package shared across layers. Each layer defines its own sentinels and returns only those, or errors with no sentinel.

Sentinels are declared in the file that returns them. There is no `errors.go` file and no wrapping helper. Each call site writes its own `fmt.Errorf` or `errors.New`.

## Storage

Storage sentinels live in a storage-owned leaf package with no internal imports, such as `internal/database/dberr`. The service can import it even when the adapter imports the domain packages.

One mapper, shared by every resource, converts driver errors. It maps a missing row to `ErrNotFound`, a unique violation to `ErrAlreadyExists`, and a foreign key violation to `ErrForeignKey`. The adapter does not return a service sentinel, and it does not import the transport layer or `net/http`.

Any other backend error is logged with its operation context, then replaced with `errors.New`. The driver error does not remain on the chain.

## Service

The service switches on the storage sentinels inline at each call. It returns its own sentinel with the entry that failed, such as `fmt.Errorf("category %q: %w", slug, ErrNotFound)`. It does not rename `ErrAlreadyExists`.

An unexpected storage error is not wrapped, so no storage error reaches transport. The service returns a new error that names the operation and entry, such as `fmt.Errorf("get category %q: internal error", slug)`. The adapter has already logged the cause. An expected failure is not logged before it is returned.

## Transport

Each call site switches on the sentinels that call can return. It does not switch on every sentinel the resource defines, and it does not call a handler shared across calls or domains. A call that cannot return a sentinel does not list it, so that sentinel cannot be written as a status for the wrong operation.

A recognized error is written with `Error()` and the status for that sentinel. Any other error from that call is a 500 written with `Error()`. The service has already replaced driver text, so the message names the operation and entry. The handler logs that error.
