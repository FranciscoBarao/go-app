package listopt

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/FranciscoBarao/catalog/internal/middleware"
)

// FieldKind is the query-time type of a list field (not the Go struct type).
type FieldKind int

const (
	// KindString allows like/eq by default.
	KindString FieldKind = iota
	// KindInt allows eq and numeric comparisons; values must parse as numbers.
	KindInt
)

const (
	Sortable      = true
	NotSortable   = false
	Filterable    = true
	NotFilterable = false
)

// FieldSpec describes one public list/query field.
type FieldSpec struct {
	Column     string
	Kind       FieldKind
	Sortable   bool
	Filterable bool
}

// Allowlist maps a public JSON field name to its query metadata.
type Allowlist map[string]FieldSpec

// StringField is a sortable and/or filterable string column.
func StringField(column string, sortable, filterable bool) FieldSpec {
	return FieldSpec{
		Column:     column,
		Kind:       KindString,
		Sortable:   sortable,
		Filterable: filterable,
	}
}

// IntField is a sortable and/or filterable integer column.
func IntField(column string, sortable, filterable bool) FieldSpec {
	return FieldSpec{
		Column:     column,
		Kind:       KindInt,
		Sortable:   sortable,
		Filterable: filterable,
	}
}

// SortOnly creates a sort-only field spec.
func SortOnly(column string) FieldSpec {
	return FieldSpec{
		Column:     column,
		Sortable:   Sortable,
		Filterable: NotFilterable,
	}
}

// ParseSort validates a sort request. Empty field and order is a no-op (zero Sort).
func (s Allowlist) ParseSort(field, order string) (Sort, error) {
	switch {
	case field == "" && order == "":
		return Sort{}, nil
	case field == "" || order == "":
		return Sort{}, middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, can't be empty")
	case order != "asc" && order != "desc":
		return Sort{}, middleware.NewError(http.StatusUnprocessableEntity, "Malformed sortBy query parameter, order should be asc or desc")
	}

	f, err := s.lookup(field)
	if err != nil {
		return Sort{}, err
	}
	if !f.Sortable {
		return Sort{}, middleware.NewError(http.StatusUnprocessableEntity, "Field not available for this operation")
	}
	return Sort{Column: f.Column, Order: order}, nil
}

// ParseFilter validates a filter request and returns a clause ready for WithFilter.
func (s Allowlist) ParseFilter(field, op, value string) (Filter, error) {
	if field == "" || value == "" {
		return Filter{}, middleware.NewError(http.StatusUnprocessableEntity, "Malformed filter, field and value can't be empty")
	}

	parsed, err := parseOperator(op)
	if err != nil {
		return Filter{}, err
	}

	f, err := s.lookup(field)
	if err != nil {
		return Filter{}, err
	}
	if !f.Filterable {
		return Filter{}, middleware.NewError(http.StatusUnprocessableEntity, "Field not available for this operation")
	}
	if !allowsOperator(f, parsed) {
		return Filter{}, middleware.NewError(http.StatusUnprocessableEntity, operatorRejectedMessage(f.Kind, parsed))
	}
	if f.Kind == KindInt {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return Filter{}, middleware.NewError(http.StatusUnprocessableEntity, "Filter value must be numeric for this operator")
		}
	}
	return Filter{Column: f.Column, Operator: parsed, Value: value, ValueKind: f.Kind}, nil
}

func (s Allowlist) lookup(field string) (FieldSpec, error) {
	f, ok := s[field]
	if !ok {
		return FieldSpec{}, middleware.NewError(http.StatusUnprocessableEntity, "No field with this name")
	}
	return f, nil
}

func parseOperator(op string) (Operator, error) {
	if op == "" || op == string(Like) {
		return Like, nil
	}
	switch Operator(op) {
	case Eq, Lt, Le, Gt, Ge:
		return Operator(op), nil
	default:
		return "", middleware.NewError(http.StatusUnprocessableEntity, "Operator not allowed")
	}
}

func (f FieldSpec) allowedOperators() []Operator {
	if f.Kind == KindString {
		return []Operator{Like, Eq}
	}
	return []Operator{Eq, Lt, Le, Gt, Ge}
}

func allowsOperator(f FieldSpec, operator Operator) bool {
	return slices.Contains(f.allowedOperators(), operator)
}

func isNumericOperator(operator Operator) bool {
	return operator == Lt || operator == Le || operator == Gt || operator == Ge
}

func operatorRejectedMessage(kind FieldKind, operator Operator) string {
	switch {
	case operator == Like:
		return "Partial match only available for string fields"
	case isNumericOperator(operator) && kind == KindString:
		return "Numeric operators not available for string fields"
	default:
		return "Operator not allowed"
	}
}
