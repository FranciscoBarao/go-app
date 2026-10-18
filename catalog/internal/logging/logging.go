package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// Init configures the global zerolog level and default context logger. Call once at startup.
// Output is zerolog's console writer. JSON requires replacing that writer; there is no format setting.
// An unknown level is debug, and Init logs a warning that names the rejected value.
func Init(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)

	// ConsoleWriter is the chosen output. JSON means replacing this writer with os.Stderr.
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &l
	if err != nil {
		l.Warn().Str("rejected", level).Msg("unknown log level")
	}
}

// FromCtx returns the logger stored in ctx, or the default context logger.
// Returns a disabled (no-op) logger if no logger is configured.
func FromCtx(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx)
	if l.GetLevel() == zerolog.Disabled && zerolog.DefaultContextLogger == nil {
		nop := zerolog.Nop()
		return &nop
	}
	return l
}

// WithLogger returns a new context with the given logger attached.
func WithLogger(ctx context.Context, l *zerolog.Logger) context.Context {
	return l.WithContext(ctx)
}
