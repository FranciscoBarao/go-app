package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

// Request-shape errors returned while decoding and validating HTTP input.
var (
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrRequestTooLarge = errors.New("request body must not be larger than 1MB")
)

// ValidateStruct executes govalidator to check if sturct fields have the correct, previously defined values
func ValidateStruct(value any) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		return fmt.Errorf("model validation failed: %s: %w", err.Error(), ErrBadRequest)
	}
	return nil
}

// GetFieldFromURL extracts a field from URL
func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}

// GetUsernameFromToken extracts the username from context
func GetUsernameFromToken(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(oauth.ClaimsContext).(map[string]string)
	if !ok {
		return "", fmt.Errorf("missing or invalid auth claims: %w", ErrUnauthorized)
	}

	username, ok := claims["username"]
	if !ok {
		return "", errors.New("username not present")
	}
	return username, nil
}
