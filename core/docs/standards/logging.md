# Logging

Service repositories log with `github.com/rs/zerolog`. The logger is created once at startup and carried on `context.Context`.

Output is zerolog's console writer. That writer is the project choice. JSON output means replacing the writer in the logging setup. There is no environment variable that selects the format.

## Level

`LOG_LEVEL` sets the global level: `trace`, `debug`, `info`, `warn`, `error`, or `fatal`. An unset value is `debug`, for local development. A deployed environment sets `LOG_LEVEL=info`. An unknown value is `debug`, and startup logs a warning with the rejected value in `rejected`.

| Level | When |
| --- | --- |
| trace | A short investigation. Off the rest of the time. |
| debug | Detail a developer needs. Off when the level is info. |
| info | Process start, process stop, and a successful change worth keeping. |
| warn | A retry, fallback, or other degraded path that still finished. |
| error | An unexpected failure. |
| fatal | The process cannot continue. Only `main` uses it, and the process exits. |

A missing record, a validation failure, and a duplicate entry are normal results. They are not error lines. An error the function returns is logged once, at the layer that handles it or replaces it.

## Fields

The message names the event and stays stable. Fields carry the values that change. The same event uses the same field names at every call site.

A line may carry a slug, an id, a status, a duration, and a request id. It does not carry credentials, tokens, cookies, personal data, raw bodies, SQL, table names, constraint names, or stored row values.

## Request scope

HTTP services install `logging.RequestMiddleware` on the router. `WithRequestIDHeader` sets the header read from the request and written on the response. The default header is `X-Request-ID`.

The middleware accepts an incoming request id when it is a UUID, and stores it in lowercase. Any other value is replaced with a new UUID. It stores a child logger on the request context with the field `request_id`, and writes one access line when the handler returns. That line is info and includes `method`, `path`, and `duration_ms`. `path` is the URL path, without the query string.

Handlers, services, and storage read the logger with `logging.FromCtx` and pass the same context into every call. A line logged with `context.Background()` has no request id.

One access line per request. A second access-log middleware is not added.

## Layers

Transport logs the unexpected failure it turns into a 500. It does not log a status it already returns for a known sentinel.

The service does not log an error it returns unchanged.

Storage logs an unexpected backend error with the caller's context, then returns a safe error. The log may include the driver code. It does not include table names, constraint names, or row values. Expected conditions (`ErrNotFound`, `ErrAlreadyExists`, `ErrForeignKey`) are not logged there.
