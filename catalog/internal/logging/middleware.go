package logging

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const defaultRequestIDHeader = "X-Request-ID"

type middlewareOptions struct {
	requestIDHeader string
}

// Option configures the request logger.
type Option func(*middlewareOptions)

// WithRequestIDHeader sets the header read from the incoming request and written on the response.
// An empty name keeps the default, X-Request-ID.
func WithRequestIDHeader(header string) Option {
	return func(options *middlewareOptions) {
		if header != "" {
			options.requestIDHeader = header
		}
	}
}

// RequestMiddleware returns HTTP middleware that attaches a request-scoped logger.
// The logger carries request_id. The middleware writes one access line when the handler returns.
func RequestMiddleware(opts ...Option) func(http.Handler) http.Handler {
	options := middlewareOptions{requestIDHeader: defaultRequestIDHeader}
	for _, opt := range opts {
		opt(&options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := requestIDFrom(r, options.requestIDHeader)
			w.Header().Set(options.requestIDHeader, requestID)

			logger := FromCtx(r.Context()).With().Str("request_id", requestID).Logger()
			started := time.Now()

			next.ServeHTTP(w, r.WithContext(logger.WithContext(r.Context())))

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int64("duration_ms", time.Since(started).Milliseconds()).
				Msg("request completed")
		})
	}
}

func requestIDFrom(r *http.Request, header string) string {
	id, err := uuid.Parse(r.Header.Get(header))
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}
