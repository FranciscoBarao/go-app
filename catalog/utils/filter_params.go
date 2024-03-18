package utils

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/middleware/logging"
)

// Examples of filters that work:
// name.a 		---> name LIKE ?    %a%
// price.le.10  ---> price <= ?     10
// name.eq.asd  ---> name == ?     asd
// GetFilters gets the filters
func GetFilters(model interface{}, filterBy string) (string, string, error) {
	if filterBy == "" {
		return "", "", nil // No filter -> No error
	}

	logging.FromCtx(context.Background()).Debug().Str("filter_by", filterBy).Msg("filtering")
	if err := validateFilter(model, filterBy); err != nil { // Validates Filters
		return "", "", err
	}

	filterBody, filterValue := getFilterBodyAndValue(filterBy) // After validating, constructs Body and Value to be used in GetAll
	return filterBody, filterValue, nil
}

// validateFilter validates if filter parameters are correct by length, emptiness and type
func validateFilter(model interface{}, filterBy string) error {
	log := logging.FromCtx(context.Background())

	var field, operator, value string

	splits := strings.Split(filterBy, ".")
	switch size := len(splits); {
	case size == 2:
		field, value = splits[0], splits[1]

	case size == 3:
		field, operator, value = splits[0], splits[1], splits[2]
		if err := validateOperator(operator); err != nil {
			return err
		}

	default:
		log.Error().Str("filter_by", filterBy).Msg("malformed query parameter, should be field.value or field.operator.value")
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, should be field.value or field.operator.value")
	}

	if field == "" || value == "" {
		log.Error().Msg("filter malformed, empty parameters")
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, can't be empty")
	}

	// Validate Field & Value
	return validateFieldAndValue(model, field, value, operator)
}

// validateFieldAndValue checks if the field exists in the struct and if the value is of the correct type
func validateFieldAndValue(model interface{}, fieldName, value, operator string) error {
	log := logging.FromCtx(context.Background())

	fields := reflect.VisibleFields(reflect.TypeOf(model)) // Get all fields of Struct
	for _, field := range fields {
		if strings.ToLower(field.Name) == fieldName { // If there is a Field with this name
			if operator == "" && field.Type.String() != "string" { // If there are only 2 field params and its not a string -> error -> E.g price.10
				log.Error().Str("field_name", fieldName).Msg("filter malformed, expected string")
				return middleware.NewError(http.StatusUnprocessableEntity, "Filter Malformed, field not a string")
			}
			if operator != "" && field.Type.String() == "string" { // If there are 3 field params and its a string -> error -> E.g name.gt.asd
				log.Error().Str("field_name", fieldName).Msg("filter malformed, expected non string")
				return middleware.NewError(http.StatusUnprocessableEntity, "Filter Malformed, field can't be a string")
			}
			return isValidType(field.Type.String(), value) // Field exists and is of the correct type
		}
	}
	log.Error().Str("field_name", fieldName).Interface("model", model).Msg("no filterable field in struct")
	return middleware.NewError(http.StatusUnprocessableEntity, "No filterable field with the provided name")
}

// isValidType receives a value and validates if it reflects the provided type.
func isValidType(typ, value string) error {
	switch typ {
	case "string":
		if isAlphanumeric(value) { // If String is alphanumeric
			return nil
		}
	case "int":
		if _, err := strconv.Atoi(value); err == nil { // If convertion to Integer was valid
			return nil
		}
	case "float64":
		if _, err := strconv.ParseFloat(value, 64); err == nil { // If convertion to Float64 was valid
			return nil
		}
	case "float32":
		if _, err := strconv.ParseFloat(value, 32); err == nil { // If convertion to Float32 was valid
			return nil
		}
	}
	logging.FromCtx(context.Background()).Error().Msg("field convertion faile due to mistype")
	return middleware.NewError(http.StatusUnprocessableEntity, "Incorrect field type")
}

// validateOperator validates the operator in the URL parameter
func validateOperator(operator string) error {
	var allowedOperators = []string{"lt", "le", "gt", "ge", "eq"}
	if !stringInSlice(operator, allowedOperators) {
		logging.FromCtx(context.Background()).Error().Str("operator", operator).Msg("unknown operator")
		return middleware.NewError(http.StatusUnprocessableEntity, "Operator not allowed")
	}
	return nil
}

// getFilterBodyAndValue gets FilterBody and Value for GetFilters
func getFilterBodyAndValue(filterBy string) (string, string) {
	splits := strings.Split(filterBy, ".")

	var field, operator, value string
	field = splits[0]

	// Must be String partial find -> all others have 3 parameters
	if len(splits) == 2 {
		value = splits[1]
		return field + " LIKE ?", "%" + value + "%"
	} else {
		operator = splits[1]
		value = splits[2]
		return field + " " + operatorToString(operator) + " ?", value
	}
}

// operatorToString converts operator language to string literal (eq -> ==)
func operatorToString(operator string) string {
	switch operator {
	case "lt":
		return "<"
	case "le":
		return "<="
	case "gt":
		return ">"
	case "ge":
		return ">="
	case "eq":
		return "=="
	default:
		return ""
	}
}
