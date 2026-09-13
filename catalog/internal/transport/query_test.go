package transport

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/boardgame"
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/stretchr/testify/require"
)

func TestNewPaginatedResponse_TotalPages(t *testing.T) {
	cases := []struct {
		name      string
		total     int
		pageSize  int
		wantPages int
	}{
		{"exact multiple", 50, 10, 5},
		{"remainder rounds up", 55, 10, 6},
		{"single partial page", 3, 10, 1},
		{"empty", 0, 10, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := newPaginatedResponse([]int{}, c.total, listopt.Pagination{Page: 1, PageSize: c.pageSize})
			require.Equal(t, c.wantPages, resp.TotalPages)
			require.Equal(t, c.total, resp.TotalItems)
		})
	}
}

func TestNewPaginatedResponse_NilDataBecomesEmptySlice(t *testing.T) {
	resp := newPaginatedResponse[int](nil, 0, listopt.Pagination{Page: 1, PageSize: 10})
	require.NotNil(t, resp.Data)
	require.Len(t, resp.Data, 0)
}

func TestQueryRequest_ToQuery(t *testing.T) {
	q := QueryRequest{
		Pagination: PaginationRequest{Page: 2, PageSize: 20},
		Sort:       &SortRequest{Field: "name", Order: "asc"},
		Filters: []FilterRequest{
			{Field: "min_players", Op: "ge", Value: "3"},
			{Field: "name", Value: "cat"}, // no op -> like
		},
	}

	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Equal(t, "name", p.Sort.Column)
	require.Equal(t, "asc", p.Sort.Order)
	require.Len(t, p.Filters, 2)
	require.Equal(t, "min_players", p.Filters[0].Column)
	require.Equal(t, listopt.Ge, p.Filters[0].Operator)
	require.Equal(t, listopt.Like, p.Filters[1].Operator)
	require.Equal(t, 2, p.Pagination.Page)
	require.Equal(t, 20, p.Pagination.PageSize)
}

func TestQueryRequest_ToQuery_FilterValuesWithDots(t *testing.T) {
	q := QueryRequest{
		Filters: []FilterRequest{
			{Field: "name", Op: "like", Value: "foo.bar"},
			{Field: "name", Op: "eq", Value: "foo.bar.baz"},
			{Field: "min_players", Op: "ge", Value: "3.5"},
		},
	}

	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Len(t, p.Filters, 3)
	require.Equal(t, listopt.Like, p.Filters[0].Operator)
	require.Equal(t, "foo.bar", p.Filters[0].Value)
	require.Equal(t, listopt.Eq, p.Filters[1].Operator)
	require.Equal(t, "foo.bar.baz", p.Filters[1].Value)
	require.Equal(t, listopt.Ge, p.Filters[2].Operator)
	require.Equal(t, "3.5", p.Filters[2].Value)
}

func TestQueryRequest_ToQuery_MinAgeNumeric(t *testing.T) {
	q := QueryRequest{Filters: []FilterRequest{{Field: "min_age", Op: "ge", Value: "8"}}}
	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Len(t, p.Filters, 1)
	require.Equal(t, "min_age", p.Filters[0].Column)
	require.Equal(t, listopt.KindInt, p.Filters[0].ValueKind)
}

func TestQueryRequest_ToQuery_OmittedFields(t *testing.T) {
	cases := []struct {
		name      string
		allowlist listopt.Allowlist
		q         QueryRequest
	}{
		{"description sort", boardgame.QueryAllowlist, QueryRequest{Sort: &SortRequest{Field: "description", Order: "asc"}}},
		{"description filter", boardgame.QueryAllowlist, QueryRequest{Filters: []FilterRequest{{Field: "description", Op: "like", Value: "x"}}}},
		{"bgg_id sort", boardgame.QueryAllowlist, QueryRequest{Sort: &SortRequest{Field: "bgg_id", Order: "asc"}}},
		{"bgg_id filter", boardgame.QueryAllowlist, QueryRequest{Filters: []FilterRequest{{Field: "bgg_id", Op: "eq", Value: "1"}}}},
		{"deleted_at sort", boardgame.QueryAllowlist, QueryRequest{Sort: &SortRequest{Field: "deleted_at", Order: "asc"}}},
		{"categories sort", boardgame.QueryAllowlist, QueryRequest{Sort: &SortRequest{Field: "categories", Order: "asc"}}},
		{"bio filter", contributor.QueryAllowlist, QueryRequest{Filters: []FilterRequest{{Field: "bio", Op: "like", Value: "x"}}}},
		{"created_at filter", boardgame.QueryAllowlist, QueryRequest{Filters: []FilterRequest{{Field: "created_at", Op: "eq", Value: "x"}}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := c.q.toQuery(c.allowlist)
			require.Error(t, err)
		})
	}
}

func TestQueryRequest_ToQuery_InvalidSortField(t *testing.T) {
	q := QueryRequest{Sort: &SortRequest{Field: "categories", Order: "asc"}}
	_, err := q.toQuery(boardgame.QueryAllowlist)
	require.Error(t, err)
}

func TestQueryRequest_ToQuery_InvalidFilterOp(t *testing.T) {
	q := QueryRequest{Filters: []FilterRequest{{Field: "name", Op: "bogus", Value: "x"}}}
	_, err := q.toQuery(boardgame.QueryAllowlist)
	require.Error(t, err)
}

func TestQueryRequest_ToQuery_DefaultsPagination(t *testing.T) {
	q := QueryRequest{} // no pagination
	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Equal(t, listopt.DefaultPage, p.Pagination.Page)
	require.Equal(t, listopt.DefaultPageSize, p.Pagination.PageSize)
}

// TestPageDocsMatchEnvelope guards the Swagger-only page types against drifting
// from the envelope the handlers actually write.
func TestPageDocsMatchEnvelope(t *testing.T) {
	cases := []struct {
		name     string
		envelope any
		doc      any
	}{
		{"boardgame", PaginatedResponse[boardgame.Boardgame]{}, BoardgamePage{}},
		{"category", PaginatedResponse[category.Category]{}, CategoryPage{}},
		{"mechanism", PaginatedResponse[mechanism.Mechanism]{}, MechanismPage{}},
		{"contributor", PaginatedResponse[contributor.Contributor]{}, ContributorPage{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			envelope, err := json.Marshal(c.envelope)
			require.NoError(t, err)
			doc, err := json.Marshal(c.doc)
			require.NoError(t, err)
			require.JSONEq(t, string(envelope), string(doc))
		})
	}
}

// parseURL builds a list request from a raw query string, as the GET list
// endpoints do.
func parseURL(t *testing.T, rawQuery string) (QueryRequest, error) {
	t.Helper()
	values, err := url.ParseQuery(rawQuery)
	require.NoError(t, err)
	return newQueryRequestFromURL(values)
}

func TestNewQueryRequestFromURL(t *testing.T) {
	q, err := parseURL(t, "page=2&pageSize=20&sort=name.asc&filter=min_players.ge.3&filter=name.like.cat&include_deleted=true")
	require.NoError(t, err)
	require.True(t, q.IncludeDeleted)

	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Equal(t, 2, p.Pagination.Page)
	require.Equal(t, 20, p.Pagination.PageSize)
	require.Equal(t, "name", p.Sort.Column)
	require.Equal(t, "asc", p.Sort.Order)
	require.Len(t, p.Filters, 2)
	require.Equal(t, "min_players", p.Filters[0].Column)
	require.Equal(t, listopt.Ge, p.Filters[0].Operator)
	require.Equal(t, "3", p.Filters[0].Value)
	require.Equal(t, "name", p.Filters[1].Column)
	require.Equal(t, listopt.Like, p.Filters[1].Operator)
	require.Equal(t, "cat", p.Filters[1].Value)
}

func TestNewQueryRequestFromURL_NoParamsUsesDefaults(t *testing.T) {
	q, err := parseURL(t, "")
	require.NoError(t, err)
	require.False(t, q.IncludeDeleted)
	require.Nil(t, q.Sort)
	require.Empty(t, q.Filters)

	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Equal(t, listopt.DefaultPage, p.Pagination.Page)
	require.Equal(t, listopt.DefaultPageSize, p.Pagination.PageSize)
	require.Empty(t, p.Sort.Column)
}

func TestNewQueryRequestFromURL_FilterValueKeepsDots(t *testing.T) {
	q, err := parseURL(t, "filter=name.like.foo.bar&filter=min_players.ge.3.5")
	require.NoError(t, err)
	require.Len(t, q.Filters, 2)
	require.Equal(t, "foo.bar", q.Filters[0].Value)
	require.Equal(t, "3.5", q.Filters[1].Value)

	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Equal(t, "foo.bar", p.Filters[0].Value)
	require.Equal(t, "3.5", p.Filters[1].Value)
}

func TestNewQueryRequestFromURL_OutOfRangePaginationIsClamped(t *testing.T) {
	q, err := parseURL(t, "page=0&pageSize=500")
	require.NoError(t, err)

	p, err := q.toQuery(boardgame.QueryAllowlist)
	require.NoError(t, err)
	require.Equal(t, listopt.DefaultPage, p.Pagination.Page)
	require.Equal(t, listopt.MaxPageSize, p.Pagination.PageSize)
}

func TestNewQueryRequestFromURL_MalformedParams(t *testing.T) {
	cases := []struct {
		name     string
		rawQuery string
	}{
		{"non-numeric page", "page=abc"},
		{"non-numeric pageSize", "pageSize=abc"},
		{"filter without operator", "filter=name.cat"},
		{"filter with only a field", "filter=name"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := parseURL(t, c.rawQuery)

			var mr *middleware.MalformedRequest
			require.ErrorAs(t, err, &mr)
			require.Equal(t, http.StatusUnprocessableEntity, mr.GetStatus())
		})
	}
}

// Sort validation happens in toQuery, so a malformed sort parameter must
// survive parsing as-is rather than being silently dropped.
func TestNewQueryRequestFromURL_MalformedSortRejectedByToQuery(t *testing.T) {
	for _, rawQuery := range []string{"sort=name", "sort=.asc", "sort=name.", "sort=name.asc.desc", "sort=name.sideways", "sort=categories.asc"} {
		t.Run(rawQuery, func(t *testing.T) {
			q, err := parseURL(t, rawQuery)
			require.NoError(t, err)
			_, err = q.toQuery(boardgame.QueryAllowlist)
			require.Error(t, err)
		})
	}
}
