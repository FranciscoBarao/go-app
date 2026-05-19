package utils

import (
	"context"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// ValidateStruct executes govalidator to check if sturct fields have the correct, previously defined values
func ValidateStruct(value any) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		middleware.FromCtx(context.Background()).Error().Err(err).Msg("model validation failed")
		return middleware.NewError(http.StatusBadRequest, "Error - Model validation failed")
	}
	return nil
}

// GetFieldFromURL extracts a field from URL
func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}

// GetUsernameFromToken extracts the username from context
func GetUsernameFromToken(r *http.Request) (string, error) {
	claims := r.Context().Value(oauth.ClaimsContext).(map[string]string)
	username, ok := claims["username"]
	if !ok {
		return "", middleware.NewError(http.StatusInternalServerError, "Error - Username not present")
	}
	return username, nil
}
