package database

import (
	"fmt"
	"strconv"
	"strings"

	dbsql "github.com/FranciscoBarao/catalog/internal/database/sql"
	"github.com/FranciscoBarao/catalog/internal/listopt"
)

// filterConditions builds AND-joined filter fragments and their args, binding
// placeholders sequentially starting at startIdx. The base query is expected to
// already contain a WHERE clause (e.g. "WHERE deleted_at IS NULL"), so each
// fragment is prefixed with " AND ".
func filterConditions(filters []listopt.Filter, startIdx int) (string, []any) {
	var query strings.Builder
	var args []any
	idx := startIdx

	for _, f := range filters {
		if f.Column == "" {
			continue
		}
		if f.Operator == listopt.Like {
			fmt.Fprintf(&query, " AND %s ILIKE $%d", f.Column, idx)
			args = append(args, "%"+f.Value+"%")
		} else {
			fmt.Fprintf(&query, " AND %s %s $%d", f.Column, operatorToSQL(f.Operator), idx)
			args = append(args, parseFilterValue(f.Value, f.Kind))
		}
		idx++
	}

	return query.String(), args
}

// buildCountQuery composes a COUNT query from a base count statement and the
// filter clauses in params.
func buildCountQuery(baseCount string, params listopt.Params) (string, []any) {
	conds, args := filterConditions(params.Filters, 1)
	return baseCount + conds, args
}

// buildPaginatedQuery composes a paginated SELECT from a base select statement,
// applying filters, optional sort, and a LIMIT/OFFSET window.
func buildPaginatedQuery(baseSelect string, params listopt.Params) (string, []any) {
	conds, args := filterConditions(params.Filters, 1)
	query := baseSelect + conds

	// Always sort so LIMIT/OFFSET pagination is deterministic. Fall back to id
	// when no sort is requested, and append id as a tiebreaker otherwise so
	// pages stay stable for non-unique sort columns.
	switch params.Sort.Column {
	case "":
		query += dbsql.OrderByIDFallback
	case "id":
		query += fmt.Sprintf(dbsql.OrderBy, params.Sort.Column, params.Sort.Order)
	default:
		query += fmt.Sprintf(dbsql.OrderBy, params.Sort.Column, params.Sort.Order)
		query += dbsql.OrderByIDTiebreak
	}

	query += fmt.Sprintf(dbsql.LimitOffset, len(args)+1, len(args)+2)
	args = append(args, params.Pagination.Limit(), params.Pagination.Offset())

	return query, args
}

func operatorToSQL(operator listopt.Operator) string {
	switch operator {
	case listopt.Lt:
		return "<"
	case listopt.Le:
		return "<="
	case listopt.Gt:
		return ">"
	case listopt.Ge:
		return ">="
	case listopt.Eq:
		return "="
	default:
		panic(fmt.Sprintf("unreachable: invalid filter operator %q", operator))
	}
}

// parseFilterValue binds a filter value with a type appropriate for the column.
// KindInt columns are coerced to int/float so pgx can encode them for numeric
// comparisons and equality. String columns keep the raw value, so an eq on a
// string column with a numeric-looking value (e.g. "123") is matched as text
// rather than silently coerced to a number.
func parseFilterValue(value string, kind listopt.FieldKind) any {
	if kind != listopt.KindInt {
		return value
	}
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return value
}
