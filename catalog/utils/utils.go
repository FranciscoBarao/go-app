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

func StringInSlice(value string, list []string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// Method that checks if a string is alphanumeric
func IsAlphanumeric(word string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(word)
}

func ValidateStruct(value interface{}) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		logging.FromCtx(context.Background()).Error().Err(err).Msg("model validation failed")
		return middleware.NewError(http.StatusForbidden, "Error - Model validation failed")
	}
	return nil
}

func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}

func GetUsernameFromToken(r *http.Request) (string, error) {
	claims := r.Context().Value(oauth.ClaimsContext).(map[string]string)

	if username, ok := claims["username"]; ok {
		return username, nil
	}

	return "", middleware.NewError(http.StatusInternalServerError, "Error - Username not present")
}
