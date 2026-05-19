package utils

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// GetFilter validates the filterBy parameter and returns the DB column, operator, and value.
//
// Supported formats:
//
//	"field.value"          → partial string match (op = OpLike)
//	"field.eq.value"       → exact equality
//	"field.lt|le|gt|ge.value" → numeric comparison
//
// Examples:
//
//	GetFilter(Boardgame{}, "name.catan")        -> ("name", OpLike, "catan", nil)
//	GetFilter(Boardgame{}, "name.eq.Catan")     -> ("name", OpEq, "Catan", nil)
//	GetFilter(Boardgame{}, "playernumber.lt.5") -> ("player_number", OpLt, "5", nil)
//	GetFilter(Boardgame{}, "")                  -> ("", "", "", nil)
func GetFilter(model any, filterBy string) (string, listopt.Op, string, error) {
	if filterBy == "" {
		return "", "", "", nil
	}

	splits := strings.Split(filterBy, ".")
	var field, value string
	var op listopt.Op

	switch len(splits) {
	case 2:
		field, op, value = splits[0], listopt.OpLike, splits[1]
	case 3:
		field, value = splits[0], splits[2]
		var err error
		op, err = validateOperator(splits[1])
		if err != nil {
			return "", "", "", err
		}
	default:
		return "", "", "", middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, should be field.value or field.operator.value")
	}

	if field == "" || value == "" {
		return "", "", "", middleware.NewError(http.StatusUnprocessableEntity, "Malformed filterBy query parameter, can't be empty")
	}

	info, err := resolveField(model, field)
	if err != nil {
		return "", "", "", err
	}

	if err := validateFilter(info.Type, op, value); err != nil {
		return "", "", "", err
	}

	return info.Column, op, value, nil
}

func validateOperator(op string) (listopt.Op, error) {
	switch op {
	case "eq":
		return listopt.OpEq, nil
	case "lt":
		return listopt.OpLt, nil
	case "le":
		return listopt.OpLe, nil
	case "gt":
		return listopt.OpGt, nil
	case "ge":
		return listopt.OpGe, nil
	default:
		return "", middleware.NewError(http.StatusUnprocessableEntity, "Operator not allowed")
	}
}

func validateFilter(fieldType string, op listopt.Op, value string) error {
	isNumericOp := op == listopt.OpLt || op == listopt.OpLe || op == listopt.OpGt || op == listopt.OpGe

	if op == listopt.OpLike && !isStringType(fieldType) {
		return middleware.NewError(http.StatusUnprocessableEntity, "Partial match only available for string fields")
	}

	if isNumericOp {
		if isStringType(fieldType) {
			return middleware.NewError(http.StatusUnprocessableEntity, "Numeric operators not available for string fields")
		}
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return middleware.NewError(http.StatusUnprocessableEntity, "Filter value must be numeric for this operator")
		}
	}

	return nil
}

func isStringType(t string) bool {
	return t == "string"
}
