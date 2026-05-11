package utils

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// GetSort validates the sortBy parameter and returns the DB column name and order.
//
// The sortBy format is "Field.Order" where Field is a case-insensitive struct field name
// and Order is "asc" or "desc". The DB column name is resolved from the field's `db` tag.
//
// Examples:
//
//	GetSort(Boardgame{}, "name.asc")          -> ("name", "asc", nil)
//	GetSort(Boardgame{}, "playernumber.desc") -> ("player_number", "desc", nil)
//	GetSort(Boardgame{}, "")                  -> ("", "", nil)
//	GetSort(Boardgame{}, "tags.asc")          -> ("", "", error: Field not sortable)
func GetSort(model any, sortBy string) (string, string, error) {
	if sortBy == "" {
		return "", "", nil
	}

	log := middleware.FromCtx(context.Background())
	log.Debug().Str("sort_by", sortBy).Msg("sorting")

	splits := strings.Split(sortBy, ".")
	if len(splits) != 2 {
		return "", "", middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, should be field.order")
	}

	field, order := splits[0], splits[1]
	if field == "" || order == "" {
		return "", "", middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, can't be empty")
	}

	if order != "asc" && order != "desc" {
		return "", "", middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, order should be asc or desc")
	}

	column, err := dbTag(model, field)
	if err != nil {
		return "", "", err
	}

	return column, order, nil
}

// dbTag finds the struct field by name (case-insensitive) and returns its db tag value.
func dbTag(model any, fieldName string) (string, error) {
	fields := reflect.VisibleFields(reflect.TypeOf(model))
	for _, f := range fields {
		if strings.EqualFold(f.Name, fieldName) {
			col := f.Tag.Get("db")
			if col == "" || col == "-" {
				return "", middleware.NewError(http.StatusUnprocessableEntity, "Field not sortable")
			}
			return col, nil
		}
	}
	return "", middleware.NewError(http.StatusUnprocessableEntity, "No field with this name")
}
