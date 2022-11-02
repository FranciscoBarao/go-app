package utils

import (
	"log"
	"marketplace/middleware"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
)

func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}

func ValidateStruct(value interface{}) error {
	if _, err := govalidator.ValidateStruct(value); err != nil {
		log.Println("Error - Model validation failed: " + err.Error())
		return middleware.NewError(http.StatusForbidden, "Error occurred, model validation failed")
	}
	return nil
}

func GetUsernameFromToken(r *http.Request) (string, error) {
	claims := r.Context().Value(oauth.ClaimsContext).(map[string]string)

	if username, ok := claims["username"]; ok {
		return username, nil
	}

	return "", middleware.NewError(http.StatusInternalServerError, "Error - Username not present")
}
