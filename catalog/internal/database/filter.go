package database

import (
	"fmt"
	"strconv"

	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
)

// filterClause builds a WHERE clause and query arg from the filter in the given listopt.Params.
// The filter value is always bound to $1.
// Returns the SQL fragment and the arg value, or empty string/nil if no filter.
func filterClause(filter listopt.Params) (string, any) {
	if filter.Filter.Column == "" {
		return "", nil
	}

	switch filter.Filter.Op {
	case listopt.OpLike:
		return fmt.Sprintf(dbsql.WhereLike, filter.Filter.Column), "%" + filter.Filter.Value + "%"
	default:
		arg := parseFilterValue(filter.Filter.Value)
		return fmt.Sprintf(dbsql.WhereOp, filter.Filter.Column, opToSQL(filter.Filter.Op)), arg
	}
}

func opToSQL(op listopt.Op) string {
	switch op {
	case listopt.OpLt:
		return "<"
	case listopt.OpLe:
		return "<="
	case listopt.OpGt:
		return ">"
	case listopt.OpGe:
		return ">="
	case listopt.OpEq:
		return "="
	default:
		panic(fmt.Sprintf("unreachable: invalid filter operator %q", op))
	}
}

func parseFilterValue(value string) any {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return value
}
