package utils

import (
	"catalog/middleware"
	"log"
	"net/http"
	"regexp"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
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
		log.Println("Error - Model validation failed: " + err.Error())
		return middleware.NewError(http.StatusForbidden, "Error occurred, model validation failed")
	}
	return nil
}

func GetFieldFromURL(r *http.Request, field string) string {
	return chi.URLParam(r, field)
}
