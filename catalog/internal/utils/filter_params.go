package utils

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// GetFilterFields validates an already-separated filter (field, operator, value)
// and returns the DB column, operator, value, and whether the column is a numeric
// type. It is used by the QUERY endpoints, where field/op/value arrive as distinct
// JSON fields, so values may legitimately contain dots (e.g. "3.5" or "foo.bar").
//
// The field is a case-insensitive struct field name resolved to its DB column
// via the model's `db` tag; fields tagged `db:"-"` are not filterable. The
// returned numeric flag reflects the column's Go kind so the database layer can
// bind the value as the correct type (e.g. an "eq" on a numeric column binds an
// int, while an "eq" on a string column binds the raw string even when it looks
// like a number, such as "123").
//
// Supported operators (op):
//
//	"" or "like"     → partial string match, ILIKE (string fields only)
//	"eq"             → exact equality (string or numeric fields)
//	"lt|le|gt|ge"    → numeric comparison (non-string fields; value must be numeric)
//
// Examples:
//
//	GetFilterFields(Boardgame{}, "name", "", "catan")         -> ("name", OpLike, "catan", false, nil)
//	GetFilterFields(Boardgame{}, "name", "like", "foo.bar")   -> ("name", OpLike, "foo.bar", false, nil)
//	GetFilterFields(Boardgame{}, "name", "eq", "123")         -> ("name", OpEq, "123", false, nil)
//	GetFilterFields(Boardgame{}, "minplayers", "eq", "3")     -> ("min_players", OpEq, "3", true, nil)
//	GetFilterFields(Boardgame{}, "minplayers", "ge", "3.5")   -> ("min_players", OpGe, "3.5", true, nil)
//	GetFilterFields(Boardgame{}, "name", "bogus", "x")        -> ("", "", "", false, error: Operator not allowed)
//	GetFilterFields(Boardgame{}, "", "eq", "x")               -> ("", "", "", false, error: field and value can't be empty)
func GetFilterFields(model any, field, op, value string) (string, listopt.Op, string, bool, error) {
	if field == "" || value == "" {
		return "", "", "", false, middleware.NewError(http.StatusUnprocessableEntity, "Malformed filter, field and value can't be empty")
	}

	var parsedOp listopt.Op
	if op == "" || op == string(listopt.OpLike) {
		parsedOp = listopt.OpLike
	} else {
		var err error
		parsedOp, err = validateOperator(op)
		if err != nil {
			return "", "", "", false, err
		}
	}

	info, err := resolveField(model, field)
	if err != nil {
		return "", "", "", false, err
	}

	if err := validateFilter(info.Type, parsedOp, value); err != nil {
		return "", "", "", false, err
	}

	return info.Column, parsedOp, value, isNumericKind(info.Kind), nil
}

// isNumericKind reports whether a reflect.Kind is an integer or floating-point
// type, i.e. a column whose filter value should be bound as a number.
func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
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
