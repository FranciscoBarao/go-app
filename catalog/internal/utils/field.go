package utils

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// fieldInfo holds resolved struct field metadata.
type fieldInfo struct {
	Column string
	Type   string
	Kind   reflect.Kind
}

// resolveField finds a struct field by name (case-insensitive) and returns its db column and Go type.
func resolveField(model any, fieldName string) (fieldInfo, error) {
	fields := reflect.VisibleFields(reflect.TypeOf(model))
	for _, f := range fields {
		if strings.EqualFold(f.Name, fieldName) {
			col := f.Tag.Get("db")
			if col == "" || col == "-" {
				return fieldInfo{}, middleware.NewError(http.StatusUnprocessableEntity, "Field not available for this operation")
			}
			return fieldInfo{Column: col, Type: f.Type.String(), Kind: f.Type.Kind()}, nil
		}
	}
	return fieldInfo{}, middleware.NewError(http.StatusUnprocessableEntity, "No field with this name")
}
