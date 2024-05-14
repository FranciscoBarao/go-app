package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

func init() {
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &l
}

// FromCtx returns the stored logger in the context or a default context logger when there is none present.
func FromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
