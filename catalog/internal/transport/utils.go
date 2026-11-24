package transport

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

var (
	// Request-shape errors returned while decoding and validating HTTP input.
	errBadRequest      = errors.New("bad request")
	errUnauthorized    = errors.New("unauthorized")
	errRequestTooLarge = errors.New("request body must not be larger than 1MB")
)

// validateRequest executes govalidator to check if sturct fields have the correct, previously defined values
func validateRequest(value any) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		return fmt.Errorf("model validation failed: %s: %w", err.Error(), errBadRequest)
	}
	return nil
}

// getFieldFromURL extracts a field from URL
func getFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}

// getUsernameFromToken extracts the username from context
func getUsernameFromToken(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(oauth.ClaimsContext).(map[string]string)
	if !ok {
		return "", fmt.Errorf("missing or invalid auth claims: %w", errUnauthorized)
	}

	username, ok := claims["username"]
	if !ok {
		return "", errors.New("username not present")
	}
	return username, nil
}
