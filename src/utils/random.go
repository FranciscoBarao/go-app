package utils

import (
	"log"
	"net/http"
	"regexp"

	"github.com/asaskevich/govalidator"
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
	_, err := govalidator.ValidateStruct(value)
	if err != nil {
		log.Println("Error - Model validation failed: " + err.Error())
		return NewError(http.StatusForbidden, "Error occurred, model validation failed")
	}
	return nil
}
