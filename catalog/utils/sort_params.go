package utils

import (
	"log"
	"net/http"
	"reflect"
	"strings"
)

// Function that checks if the sort parameters are valid for use
func validateSortParameters(model interface{}, sortBy string) error {

	splits := strings.Split(sortBy, ".")

	if len(splits) != 2 { // Validate if there are only 2 parameters
		return NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, should be field.order")
	}

	field := splits[0]
	order := splits[1]

	if field == "" || order == "" { // Validate if there are no empty parameters
		log.Printf("Error - Filter malformed, empty parameters")
		return NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, can't be empty")
	}

	if order != "desc" && order != "asc" { // Validate if order is valid
		return NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, order should be asc or desc")
	}

	err := validateField(model, field) // Validate if field exists
	if err != nil {
		return err
	}

	return nil
}

// Function that checks if the field exists in the struct
func validateField(model interface{}, fieldName string) error {

	fields := reflect.VisibleFields(reflect.TypeOf(model)) // Get all fields of Struct
	for _, field := range fields {
		if strings.ToLower(field.Name) == fieldName { // If there is a Field with this name

			err := isTypeSortable(field.Type.String()) // Checks if field is sortable
			if err != nil {
				return err
			}
			return nil // Field exists
		}
	}
	log.Printf("Error - No field in struct %v with name %s", model, fieldName)
	return NewError(http.StatusUnprocessableEntity, "No field with this name")
}

// Function that verifies if the field is sortable (E.g We cant sort by Tags)
func isTypeSortable(typ string) error {
	switch typ {
	case "string", "int", "float64", "float32":
		return nil
	}
	log.Printf("Error - Field of type %s is not sortable", typ)
	return NewError(http.StatusUnprocessableEntity, "Field not sortable")
}

// Function that constructs sort query
func constructSort(sortBy string) string {
	splits := strings.Split(sortBy, ".")
	field := splits[0]
	order := splits[1]

	return field + " " + order

}

// Main function of constructing the Sort
func GetSort(model interface{}, sortBy string) (string, error) {

	if sortBy != "" {
		log.Println("Sorting using %s " + sortBy)
		err := validateSortParameters(model, sortBy) // Validates Sort -> By length, emptiness and by order and field existence
		if err != nil {
			return "", err
		}

		sort := constructSort(sortBy) // After validating, constructs sort to be used in GetAll
		return sort, nil
	}
	return "", nil // No sort -> No error

	// Examples of sorts that work:
	// name.asc      --->  ordered by name in alphabetical ascending order
	// price.desc    --->  ordered by price in numerical descending order
}
