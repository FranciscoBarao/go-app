package utils

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/middleware/logging"
)

// GetSort constructs the whole sort
func GetSort(model interface{}, sortBy string) (string, error) {
	log := logging.FromCtx(context.Background())
	if sortBy == "" {
		return "", nil // No sort -> No error
	}

	log.Debug().Str("sort_by", sortBy).Msg("sorting")
	if err := validateSort(model, sortBy); err != nil {
		return "", err
	}

	return constructSort(sortBy), nil
}

// validateSort checks if the sort parameters are valid for use by length, emptiness, order and field existence
func validateSort(model interface{}, sortBy string) error {
	splits := strings.Split(sortBy, ".")

	if len(splits) != 2 { // Validate if there are only 2 parameters
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, should be field.order")
	}

	field, order := splits[0], splits[1]
	if field == "" || order == "" { // Validate if there are no empty parameters
		logging.FromCtx(context.Background()).Error().Msg("sort malformed with empty parameters")
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, can't be empty")
	}

	if order != "desc" && order != "asc" { // Validate if order is valid
		logging.FromCtx(context.Background()).Error().Str("order", order).Msg("sort malformed with incorrect parameters")
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, order should be asc or desc")
	}

	return validateField(model, field) // Validate if field exists
}

// validateField checks if a field exists in the struct
func validateField(model interface{}, fieldName string) error {
	fields := reflect.VisibleFields(reflect.TypeOf(model)) // Get all fields of Struct
	for _, field := range fields {
		if strings.ToLower(field.Name) == fieldName { // If there is a Field with this name
			return isTypeSortable(field.Type.String()) // Checks if field is sortable
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
	default:
		logging.FromCtx(context.Background()).Error().Str("type", typ).Msg("field is not sortable")
		return middleware.NewError(http.StatusUnprocessableEntity, "Field not sortable")
	}
}

// constructSort constructs the sort query
func constructSort(sortBy string) string {
	splits := strings.Split(sortBy, ".")
	field, order := splits[0], splits[1]
	return field + " " + order
}
