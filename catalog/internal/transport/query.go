package transport

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/utils"
)

// SortRequest describes the sort clause of a list request.
type SortRequest struct {
	Field string `json:"field"`
	Order string `json:"order"` // "asc" or "desc"
}

// FilterRequest describes a single filter clause of a list request.
type FilterRequest struct {
	Field string `json:"field"`
	Op    string `json:"op"` // like, eq, lt, le, gt, ge (empty -> like)
	Value string `json:"value"`
}

// PaginationRequest describes the requested pagination window.
type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// QueryRequest is a list request: the JSON body accepted by the QUERY endpoints,
// and the parsed form of the GET list query parameters.
type QueryRequest struct {
	Pagination     PaginationRequest `json:"pagination"`
	Sort           *SortRequest      `json:"sort,omitempty"`
	Filters        []FilterRequest   `json:"filters,omitempty"`
	IncludeDeleted bool              `json:"include_deleted,omitempty"`
}

// Query parameter names accepted by the GET list endpoints.
const (
	paramPage           = "page"
	paramPageSize       = "pageSize"
	paramSort           = "sort"
	paramFilter         = "filter"
	paramIncludeDeleted = "include_deleted"
)

// newQueryRequestFromURL builds a list request from GET query parameters, so the
// GET and QUERY entry points share the validation in toOptions.
//
//	page=2&pageSize=20      pagination window (out-of-range values are clamped)
//	sort=name.asc           field.order
//	filter=minplayers.ge.3  field.op.value, repeatable, AND-combined
//	include_deleted=true    boardgames only
//
// The operator is mandatory in filter, unlike the JSON form where it may be
// omitted to mean "like": splitting on the first two dots is what lets a value
// contain dots (e.g. filter=name.like.foo.bar).
func newQueryRequestFromURL(values url.Values) (QueryRequest, error) {
	page, err := parseIntParam(values, paramPage)
	if err != nil {
		return QueryRequest{}, err
	}
	pageSize, err := parseIntParam(values, paramPageSize)
	if err != nil {
		return QueryRequest{}, err
	}

	q := QueryRequest{
		Pagination:     PaginationRequest{Page: page, PageSize: pageSize},
		IncludeDeleted: values.Get(paramIncludeDeleted) == "true",
	}

	if sortBy := values.Get(paramSort); sortBy != "" {
		field, order, _ := strings.Cut(sortBy, ".")
		q.Sort = &SortRequest{Field: field, Order: order}
	}

	for _, raw := range values[paramFilter] {
		parts := strings.SplitN(raw, ".", 3)
		if len(parts) != 3 {
			return QueryRequest{}, middleware.NewError(http.StatusUnprocessableEntity, "Malformed filter query parameter, should be field.op.value")
		}
		q.Filters = append(q.Filters, FilterRequest{Field: parts[0], Op: parts[1], Value: parts[2]})
	}

	return q, nil
}

// parseIntParam reads an optional integer query parameter, returning 0 when it is
// absent so listopt applies its default. Out-of-range values are left for
// listopt to clamp; only unparseable ones are rejected.
func parseIntParam(values url.Values, name string) (int, error) {
	raw := values.Get(name)
	if raw == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, middleware.NewError(http.StatusUnprocessableEntity, "Malformed "+name+" query parameter, must be an integer")
	}
	return n, nil
}

// sortParam renders the requested sort in the dotted "field.order" form taken by
// utils.GetSort, or "" when no sort was requested.
func (q *QueryRequest) sortParam() string {
	if q.Sort == nil || (q.Sort.Field == "" && q.Sort.Order == "") {
		return ""
	}
	return q.Sort.Field + "." + q.Sort.Order
}

// toOptions validates the list request against the given model and converts it
// into listopt options. Field/operator validation reuses the existing sort and
// filter validators.
func (q *QueryRequest) toOptions(model any) ([]listopt.Option, error) {
	var opts []listopt.Option

	if sortBy := q.sortParam(); sortBy != "" {
		col, order, err := utils.GetSort(model, sortBy)
		if err != nil {
			return nil, err
		}
		if col != "" {
			opts = append(opts, listopt.WithSort(col, order))
		}
	}

	for _, f := range q.Filters {
		col, op, val, numeric, err := utils.GetFilterFields(model, f.Field, f.Op, f.Value)
		if err != nil {
			return nil, err
		}
		if col != "" {
			opts = append(opts, listopt.WithFilter(col, op, val, numeric))
		}
	}

	opts = append(opts, listopt.WithPagination(q.Pagination.Page, q.Pagination.PageSize))
	return opts, nil
}

// PageMeta describes the page window and result counts of a list response.
type PageMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// PaginatedResponse is the envelope returned by the GET and QUERY list endpoints.
type PaginatedResponse[T any] struct {
	Data []T `json:"data"`
	PageMeta
}

// Per-resource mirrors of PaginatedResponse, used only in the Swagger
// annotations of the GET list endpoints: swag v1.8.x renders a generic type as an
// untyped object, which would hide the item schema. They are structurally
// verified against the real envelope by TestPageDocsMatchEnvelope.
// TODO: replace with PaginatedResponse[T] in the annotations once swag is
// upgraded to a version that supports generics.
type (
	BoardgamePage struct {
		Data []boardgame.Boardgame `json:"data"`
		PageMeta
	}
	CategoryPage struct {
		Data []category.Category `json:"data"`
		PageMeta
	}
	MechanismPage struct {
		Data []mechanism.Mechanism `json:"data"`
		PageMeta
	}
	ContributorPage struct {
		Data []contributor.Contributor `json:"data"`
		PageMeta
	}
)

// newPaginatedResponse builds an envelope from results and the total item count,
// computing totalPages from the (already normalized) pagination window.
func newPaginatedResponse[T any](data []T, total int, p listopt.Pagination) PaginatedResponse[T] {
	if data == nil {
		data = []T{}
	}
	totalPages := 0
	if p.PageSize > 0 {
		totalPages = (total + p.PageSize - 1) / p.PageSize
	}
	return PaginatedResponse[T]{
		Data: data,
		PageMeta: PageMeta{
			Page:       p.Page,
			PageSize:   p.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}
