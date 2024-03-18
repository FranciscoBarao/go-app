package utils

import (
	"context"
	"net/http"
	"regexp"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/middleware/logging"
)

// ValidateStruct executes govalidator to check if sturct fields have the correct, previously defined values
func ValidateStruct(value interface{}) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		logging.FromCtx(context.Background()).Error().Err(err).Msg("model validation failed")
		return middleware.NewError(http.StatusForbidden, "Error - Model validation failed")
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

// stringInSlice checks if a specific string exists in a slice of strings
func stringInSlice(value string, list []string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// isAlphanumeric checks if a string is alphanumeric
func isAlphanumeric(word string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(word)
}
