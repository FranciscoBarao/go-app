# Config loading

A service reads its configuration once, at startup, into one struct.

Config loading via `github.com/caarlos0/env/v11`. A service uses that module when its configuration is next changed.

## Source

The process environment is the runtime source. A local `.env` file is applied by Compose or by the shell. The binary does not open that file.

A tag names each variable. The name is written in the tag. The loader does not derive it from the field name.

```go
type Config struct {
    LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
    Port     int    `env:"PORT,required"`
    OAuthKey string `env:"OAUTH_KEY,required,notEmpty"`
    Postgres PostgresConfig
}
```

`required` fails when the variable is unset. `notEmpty` fails when the value is blank. `envDefault` fills an optional field when the variable is unset or empty. A secret or a host has no default.

A nested block carries the variable name on each leaf. A shared prefix uses `envPrefix` on the parent. A repository path that has one usual location takes `envDefault`, so `Load` fills it and an environment variable can still replace it. Catalog's migration directory defaults to `internal/store/postgres/migrations`.

## Load

`Load` returns `*Config`. It calls `env.Parse` once. The returned error names every missing required variable and every value that failed to parse. The process exits before it listens.

A range or an enum that the tag cannot express is checked in Go after the parse. Catalog rejects a port outside `1-65535` there.

`LOG_LEVEL` is a field on `Config`. An unset or empty value is `debug`. An unknown value is left as read. `logging.SetLevel` warns, with the rejected value in `rejected`, and uses `debug`.

## What stays out

The config struct has no `String` method. A password, a token, or a URL with userinfo must not be formatted into a log line.

The process does not reload a listen address, a database setting, or a credential. A change to those values is a new process. A log level is applied by `logging.SetLevel` after `Load`.

## Tests

A test of `Load` passes an environment map into the parser. A test of later behavior builds the struct directly. Store tests construct `PostgresConfig` and do not read the process environment.
