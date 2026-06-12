package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/FranciscoBarao/user-management/internal/logging"
)

// ErrorHandler writes an appropriate HTTP error response based on the error type.
func ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		if mr, ok := errors.AsType[*MalformedRequest](err); ok {
			logging.FromCtx(context.Background()).Error().Int("status", mr.GetStatus()).Str("error", mr.Error()).Msg("request error")
			http.Error(w, mr.GetMessage(), mr.GetStatus())
		} else {
			logging.FromCtx(context.Background()).Error().Err(err).Msg("internal server error")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
