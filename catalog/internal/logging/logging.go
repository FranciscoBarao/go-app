package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// Init creates the process logger. Call once at startup, before configuration is read.
// The level is debug until SetLevel. Output is zerolog's console writer.
// JSON requires replacing that writer; there is no format setting.
func Init() {
	// ConsoleWriter is the chosen output. JSON means replacing this writer with os.Stderr.
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &l
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// SetLevel sets the global level from configuration.
// An unknown level is debug, and SetLevel logs a warning that names the rejected value.
func SetLevel(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		FromCtx(context.Background()).Warn().Str("rejected", level).Msg("unknown log level")
		return
	}
	zerolog.SetGlobalLevel(lvl)
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
