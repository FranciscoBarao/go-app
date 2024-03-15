package logging

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	l := log.With().Logger()
	zerolog.DefaultContextLogger = &l
}

// FromCtx returns the stored logger in the context or a default context logger when there is none present.
func FromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
