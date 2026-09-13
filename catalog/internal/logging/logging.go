package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// Init configures the global zerolog level and default context logger. Call once at startup.
func Init(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)

	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &l
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
