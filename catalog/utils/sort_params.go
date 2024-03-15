package utils

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/middleware/logging"
)

// validateSortParameters checks if the sort parameters are valid for use
func validateSortParameters(model interface{}, sortBy string) error {
	splits := strings.Split(sortBy, ".")

	if len(splits) != 2 { // Validate if there are only 2 parameters
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, should be field.order")
	}

	field := splits[0]
	order := splits[1]

	if field == "" || order == "" { // Validate if there are no empty parameters
		logging.FromCtx(context.Background()).Error().Msg("filter malformed with empty parameters")
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, can't be empty")
	}

	if order != "desc" && order != "asc" { // Validate if order is valid
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, order should be asc or desc")
	}

	return validateField(model, field) // Validate if field exists
}

// validateField checks if the field exists in the struct
func validateField(model interface{}, fieldName string) error {
	fields := reflect.VisibleFields(reflect.TypeOf(model)) // Get all fields of Struct
	for _, field := range fields {
		if strings.ToLower(field.Name) == fieldName { // If there is a Field with this name

			return isTypeSortable(field.Type.String()) // Checks if field is sortable.
		}
	}
	logging.FromCtx(context.Background()).Error().Interface("model", model).Str("field_name", fieldName).Msg("unknown field in struct")
	return middleware.NewError(http.StatusUnprocessableEntity, "No field with this name")
}

// isTypeSortable verifies if the field is sortable (E.g We cant sort by Tags)
func isTypeSortable(typ string) error {
	switch typ {
	case "string", "int", "float64", "float32":
		return nil
	}
	logging.FromCtx(context.Background()).Error().Str("type", typ).Msg("field is not sortable")
	return middleware.NewError(http.StatusUnprocessableEntity, "Field not sortable")
}

// constructSort constructs sort query
func constructSort(sortBy string) string {
	splits := strings.Split(sortBy, ".")
	field := splits[0]
	order := splits[1]

	return field + " " + order
}

// GetSort constructs the whole sort
func GetSort(model interface{}, sortBy string) (string, error) {
	log := logging.FromCtx(context.Background())
	if sortBy != "" {
		log.Debug().Str("sort_by", sortBy).Msg("sorting..")
		err := validateSortParameters(model, sortBy) // Validates Sort -> By length, emptiness, order and field existence
		if err != nil {
			return "", err
		}

		sort := constructSort(sortBy) // After validating, constructs sort to be used in GetAll
		return sort, nil
	}
	return "", nil // No sort -> No error
}
