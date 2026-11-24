# Config loading

Config is the set of values a process needs before it can work: a port, a database host, a log level, a key. Loading reads those values, checks them, and stores them on one typed object. The object is built once, at startup, and passed into the program. A handler does not read the environment on the request path.


## Where the values live

A value lives in one of three places. The process can also combine them. Later in this note, "source" means one of these places.

| Source | What it is | When it fits |
| --- | --- | --- |
| Environment | A name and a string the platform set on the process | Containers, CI, and any host that injects config from outside the image |
| File | YAML, JSON, TOML, or a `.env` file the process opens | A large tree of defaults a reviewer can read, or a secret mounted as a file |
| Flags | Arguments on the command line | A tool a person runs by hand |


### Environment variables

The platform sets the environment before the process starts. Docker Compose, Kubernetes, and a shell can all do that. The image stays the same from a laptop to production. Only the environment changes. A secret can stay out of the git tree.

The cost is the shape. Every value is a string. Nested settings become long names (`DATABASE_HOST`, `DATABASE_PORT`). A missing name is easy to miss until the process starts. A local developer needs a way to set the same names, which is usually a `.env` file applied by the shell or by Compose, not by the binary.

An empty value and an unset name are different. Unset means "use the default, or fail if the field is required". Empty means someone set the name and left it blank. A loader that treats those as the same will hide a bad deploy.


### Files

A file can hold a nested tree, comments, and defaults that a reviewer can diff. YAML, JSON, and TOML are the usual formats. The process opens a path, parses the file, and maps it onto the config object.

The cost is the path. The binary has to find the file, and the file has to be present in every environment. A secret committed next to the code is a leaked secret. A file baked into the image cannot change per environment without a rebuild.

A `.env` file is a special case. It is a list of `NAME=value` lines. Two different things are called "loading" it:

- The platform reads the file and sets the process environment. Compose `env_file` works this way. The binary only sees variables.
- The binary opens the file and parses it. That is a second source inside the process, with its own path and its own errors.

The first form keeps one code path. The second form is useful for a program that must run with no platform around it.


### Flags

A flag is a value on the command that started the process. It suits a tool a person invokes once. A long-running service rarely wants its database password on the command line, because the command is visible to other users on the host. Flags still belong in the precedence list, because a tool and a service often share one config type.


## How the values land on a type

The source is a bag of strings. The program wants a struct (or a class) with an `int` port, a duration, and a nested database block. Two ways do that mapping.


### Field by field

The code looks up each name, parses it, and assigns the field.

```go
port, ok := os.LookupEnv("PORT")
if !ok {
    return fmt.Errorf("missing environment variables: PORT")
}
```

This is easy to read for a handful of strings. Each new field is more code: the lookup, the parse, the default, and the name in the error. A nested block copies that block. The type of the field and the name of the variable can drift apart, because nothing ties them together except the author.


### Tags on the type

The struct is the schema. A tag on each field names the variable, the default, and whether the field is required. A loader walks the struct with reflection and fills it.

```go
type Config struct {
    LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
    Port     int    `env:"PORT,required"`
    Database struct {
        Host string `env:"HOST,required"`
    } `envPrefix:"DATABASE_"`
}
```

`Port` arrives as an `int`. `LogLevel` is `debug` when `LOG_LEVEL` is unset. `Database.Host` reads `DATABASE_HOST`. Adding a field is a field and a tag.

The cost is reflection and the tag language. A typo in a tag is a runtime error, not a compile error. The tag set has to be learned. A loader that guesses the variable name from the field name hides a rename: renaming `Host` to `Hostname` silently looks up a new variable. Explicit names in the tag avoid that.

The same idea exists for files. A `yaml:"host"` tag, or a JSON name, maps a file key onto the field. One struct can carry both a file tag and an env tag when both sources are in use.


## Nesting

Flat names are enough for a small service. A nested struct earns its place when one block is repeated or passed around on its own, such as a database config shared by the pool and the migrator.

Environment variables stay flat, so nesting uses a prefix. The parent adds `DATABASE_`, and the child field adds `HOST`. The variable the operator sets is still `DATABASE_HOST`.

A file nests for real:

```yaml
database:
  host: catalog-db
  port: 5432
```

The struct follows that tree. An override from the environment then targets one leaf, usually with the prefix form above.

Deep trees get expensive. A third or fourth level produces names like `CACHE_REDIS_POOL_MAX_IDLE`, and an override rule has to say which level wins. Prefer a shallow struct: one block per external system. A field with one usual location, such as a migration directory inside the repo, takes a default. An environment variable can still replace that default.


## More than one source

A real program often has a default, a file of shared settings, and an environment that overrides one or two of them. The loader applies sources in a fixed order. The later source wins.

A useful order is:

1. Defaults compiled into the struct.
2. A config file, if the process opens one.
3. The environment.
4. Flags, if the program has them.

The operator's environment beats the file, which beats the default. That lets a file hold the common case and a variable hold the one value that changes per deploy.

Document the order next to `Load`. A reader should be able to say which value won without running the process. Apply the whole order in one function, and return every missing required name in that one error. Failing on the first missing name turns a deploy into a loop of restarts.

Defaults belong on the optional fields only. A required password with a default is a password that ships inside the binary.


## Mutation before the value is stored

A mutation changes the raw string before the loader parses it onto the field. The field then stores the result, not the original text.

Typical mutations:

- Trim space and quotes that a file or a shell added.
- Expand a reference such as `postgres://${DB_HOST}/app` into a finished URL.
- Treat the value as a path, read the file, and store the contents. That is how a mounted secret (`/run/secrets/db`) becomes a password field.
- Rewrite a case, so `INFO` and `info` land as the same log level.

The hook sees the string and returns a string. It does not know the field type yet. Parsing (`int`, `time.Duration`, a custom type) happens after the hook.

Keep the hook mechanical. A rule about the business belongs after load, in ordinary code: the log level is one of a known set, the port is inside `1..65535`, the host is not empty. Those checks read clearly in a test, and they can name the field in the error. A mutation that also validates mixes two jobs and is harder to see.

Two built-in forms show up often, and they are mutations even when the API does not call them that:

- An `expand` option on a tag.
- A `file` option that means "this variable holds a path".

A general mutator is the same idea with the function supplied by the program. Use it when the built-in forms do not cover the transform. Do not use it to set a field the environment does not own.


## Reloading

Reload builds the config object again after the process is already running. Three mechanisms exist:

- A file watcher fires when the config file changes.
- The program re-reads the environment on a timer or on a signal such as `SIGHUP`.
- A tag marks a few fields as refreshable, and the program calls an update function before it uses them.

Reload fits a value the process can apply without rebuilding its world: a log level, a feature flag, a timeout. The call site has to notice the new value. A watcher that replaces a struct another goroutine is reading will race unless one lock covers both the swap and the read.

Reload is a poor fit for a value that built a long-lived resource. A new database password does nothing until the pool reconnects. A new listen port does nothing until the server rebinds. Those changes are a restart. A restart is simpler than a half-applied reload, and it fails in the open if the new value is bad.

If a service does reload, limit it to the fields that can change live, and keep the startup check on the rest. A failed reload leaves the previous object in place and logs the failure. It does not exit the process, and it does not serve with a half-filled config.


## Practices

- Build one config value at startup. Pass it down. Do not keep a package-level variable that later code mutates.
- Fail before the process listens. A missing database host is a startup error, not an error on the first request.
- Return every missing required name in one error.
- Treat unset and empty as different. Required means the name is present and, where it matters, not blank.
- Put defaults on optional fields. Leave secrets and hosts required.
- Prefer tags over a hand-written lookup once the set is more than a few fields, or as soon as a field is not a string.
- Write the variable name in the tag. Do not derive it from the field name.
- Keep the struct shallow. Use a prefix for a nested block. Give a repository path a default when it has one usual location.
- If a file and the environment both apply, write the order down: file, then environment, then flags. The environment wins over the file.
- Let the platform apply a local `.env`. The binary reads variables.
- Mutate only to clean or resolve the raw string. Check enums and ranges after load.
- Do not log the config object. A password, a token, or a full URL with a userinfo part will end up in the log. Do not give the struct a `String` method that prints those fields.
- Reload a log level or a flag if the process can apply it. Restart for a listener address, a pool, or a credential.
- In a test, build the struct directly when the test is not about loading. When the test is about loading, set the environment for that test only.
