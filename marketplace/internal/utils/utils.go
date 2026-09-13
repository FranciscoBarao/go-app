package utils

import (
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/marketplace/internal/logging"
	"github.com/FranciscoBarao/marketplace/internal/middleware"
)

// ValidateStruct executes govalidator to check struct fields.
func ValidateStruct(value any) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		logging.FromCtx(context.Background()).Error().Err(err).Msg("model validation failed")
		return middleware.NewError(http.StatusBadRequest, "Error - Model validation failed")
	}
	return nil
}

// GetFieldFromURL extracts a field from URL.
func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}

// GetUsernameFromToken extracts the username from the oauth context.
func GetUsernameFromToken(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(oauth.ClaimsContext).(map[string]string)
	if !ok {
		return "", middleware.NewError(http.StatusUnauthorized, "Error - Missing or invalid auth claims")
	}

	username, ok := claims["username"]
	if !ok {
		return "", middleware.NewError(http.StatusInternalServerError, "Error - Username not present")
	}
	return username, nil
}
