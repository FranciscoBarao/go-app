package utils

import (
	"log"
	"marketplace/middleware"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
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
