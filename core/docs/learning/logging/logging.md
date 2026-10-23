# Logging

A log records something that happened so a person or a machine can find it later. A useful line names one event, carries a few stable fields, and stays quiet the rest of the time.


## Levels

| Level | When                                                                                         |
| ----- | -------------------------------------------------------------------------------------------- |
| trace | Very fine steps inside a function. Off outside a short investigation.                        |
| debug | Diagnostic detail for a developer. Off in production.                                        |
| info  | Lifecycle and successful changes worth keeping: process start, a write that completed.       |
| warn  | The operation finished on a degraded path: a retry that succeeded, a fallback that was used. |
| error | The operation failed and someone may need to act.                                            |
| fatal | The process cannot continue. The process exits after the line.                               |

A missing record returned to the client is a normal result. Error means the system failed.

```
level=info msg="server starting" port=8080
level=error msg="create boardgame" slug=catan err="database operation failed"
```


## What to log

Log the operation, a stable identifier (slug, id), the outcome, and a duration when the call is slow or leaves the process. One event per outcome. A reader should be able to say what happened without opening the source.

```
level=info msg="boardgame created" slug=catan duration_ms=12
```


## What to leave out 

Leave out results the API already returns to the caller: a missing row, a validation failure, a duplicate slug. Those are control flow.

Leave out: 
- Secrets, tokens, passwords and personal data. 
- Session cookies 
- Raw request bodies
- SQL that includes bound values
- Errors the function also returns


## Cost

Building a line costs CPU and allocations. Watch for these:

- Formatting strings, or walking a struct with reflection, before the level check. A disabled debug line should do no work.
- Encoding a large object to JSON on the request path. Log an id and a count. Leave the payload out.
- Writing every success to disk or the network. Sample high-volume success paths. Keep every error.
- A line inside a tight loop. Log one count when the loop finishes.
- A console or pretty printer in a collected environment. That writer is slower than a JSON handler. Console output suits a laptop.


## Relevance

Each line answers who, what, and the outcome, and carries an id that groups every line from one request. Drop a line that repeats a fact already logged. Drop a field that changes on every call and does not help a search, such as a full payload or generated SQL.

The request id is the field that ties the handler, the service, and the store together. Without it, three lines from one call look like three unrelated events.


## Text

One event is one line. The message names the event and stays stable, so a search for that message still works next month. Fields carry the values that change. The message is a label, not a sentence a caller has to parse.

```
msg="create boardgame" slug=catan
```

A message that embeds the slug (`failed to create boardgame catan: pq: ...`) cannot be grouped, and it can store data that should stay out of the log.


## Structured logging

Structured logging writes key/value fields instead of one sentence. A collector filters `level=error` and `slug=catan` without a regular expression. The same event uses the same message and the same field names at every call site. `slug` in one handler and `boardgame_slug` in another splits the search.

Unstructured:

```
failed to create boardgame catan: duplicate slug
```

Structured:

```json
{"level":"error","msg":"create boardgame","slug":"catan","err":"already exists"}
```