package utils

import (
	"catalog/middleware"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Function that validates if Filter Parameters are Correct
func validateFilterParameters(model interface{}, filterBy string) error {
	splits := strings.Split(filterBy, ".")

	if len(splits) < 2 || len(splits) > 3 { // can be 2 or 3 fields
		log.Printf("Error - Malformed filterBy query parameter, should be field.value or field.operator.value")
		return middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, should be field.value or field.operator.value")
	}

	// Get field, operator, value
	var field, operator, value string
	field = splits[0]

	if len(splits) == 2 {
		if splits[0] == "" || splits[1] == "" {
			log.Printf("Error - Filter malformed, empty parameters")
			return middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, can't be empty")
		}
		value = splits[1]
	}
	if len(splits) == 3 {
		if splits[0] == "" || splits[1] == "" || splits[2] == "" {
			log.Printf("Error - Filter malformed, empty parameters")
			return middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, can't be empty")
		}

		operator = splits[1]
		// Validate Operator
		err := validateOperator(operator)
		if err != nil {
			return err
		}
		value = splits[2]
	}

	// Validate Field & Value
	return validateFieldAndValueType(model, field, value, operator)
}

// Function that checks if the field exists in the struct and if the value is of the correct type
func validateFieldAndValueType(model interface{}, fieldName, value, operator string) error {

	fields := reflect.VisibleFields(reflect.TypeOf(model)) // Get all fields of Struct
	for _, field := range fields {
		if strings.ToLower(field.Name) == fieldName { // If there is a Field with this name

			if operator == "" && field.Type.String() != "string" { // If there are only 2 field params and its not a string -> error -> E.g price.10
				log.Printf("Error - Filter malformed, %s had to be a String", fieldName)
				return middleware.NewError(http.StatusUnprocessableEntity, "Filter Malformed, field not a string")
			}
			if operator != "" && field.Type.String() == "string" { // If there are 3 field params and its a string -> error -> E.g name.gt.asd
				log.Printf("Error - Filter malformed, %s can't be a String", fieldName)
				return middleware.NewError(http.StatusUnprocessableEntity, "Filter Malformed, field can't be a string")
			}

			return isValidType(field.Type.String(), value) // Field exists and is of the correct type
		}
	}
	log.Printf("Error - No filterable field in struct %v with name %s", model, fieldName)
	return middleware.NewError(http.StatusUnprocessableEntity, "No filterable field with this name")
}

// Function that receives a value and validates if it is of the provided type.
func isValidType(typ, value string) error {
	switch typ {
	case "string":
		if IsAlphanumeric(value) { // If String is alphanumeric
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
	log.Printf("Error - Field convertion faile due to Mistype")
	return middleware.NewError(http.StatusUnprocessableEntity, "Field not of the correct type")
}

// Function that validates the operator in the URL parameter
func validateOperator(operator string) error {

	var allowedOperators = []string{"lt", "le", "gt", "ge", "eq"}
	if StringInSlice(operator, allowedOperators) {
		return nil
	}
	log.Printf("Error - No such Operator: %s", operator)
	return middleware.NewError(http.StatusUnprocessableEntity, "Operator not allowed")
}

// Function that gets FilterBody and Value for GetFilters
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

// Converts operator language to string literal (eq -> ==)
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

// Main function of Getting the Filters
func GetFilters(model interface{}, filterBy string) (string, string, error) {

	if filterBy != "" {
		log.Println("Filtering using %s " + filterBy)
		err := validateFilterParameters(model, filterBy) // Validates Filters -> By length, emptiness and type
		if err != nil {
			return "", "", err
		}

		filterBody, filterValue := getFilterBodyAndValue(filterBy) // After validating, constructs Body and Value to be used in GetAll
		return filterBody, filterValue, nil
	}
	return "", "", nil // No filter -> No error

	// Examples of filters that work:
	// name.a 		---> name LIKE ?    %a%
	// price.le.10  ---> price <= ?     10
	// name.eq.asd  ---> name == ?     asd
}
