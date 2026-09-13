package database

import (
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/stretchr/testify/require"
)

func TestFilterConditions(t *testing.T) {
	tests := []struct {
		name     string
		filters  []listopt.Filter
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "no filters",
			filters:  nil,
			wantSQL:  "",
			wantArgs: nil,
		},
		{
			name:     "like op produces ILIKE clause",
			filters:  []listopt.Filter{{Column: "name", Operator: listopt.Like, Value: "catan"}},
			wantSQL:  " AND name ILIKE $1",
			wantArgs: []any{"%catan%"},
		},
		{
			name:     "eq op on string column keeps string value",
			filters:  []listopt.Filter{{Column: "name", Operator: listopt.Eq, Value: "Catan"}},
			wantSQL:  " AND name = $1",
			wantArgs: []any{"Catan"},
		},
		{
			name:     "eq op on string column with numeric-looking value stays string",
			filters:  []listopt.Filter{{Column: "name", Operator: listopt.Eq, Value: "123"}},
			wantSQL:  " AND name = $1",
			wantArgs: []any{"123"},
		},
		{
			name:     "eq op on numeric column coerces to int",
			filters:  []listopt.Filter{{Column: "player_number", Operator: listopt.Eq, Value: "5", Kind: listopt.KindInt}},
			wantSQL:  " AND player_number = $1",
			wantArgs: []any{5},
		},
		{
			name:     "lt op with numeric value",
			filters:  []listopt.Filter{{Column: "player_number", Operator: listopt.Lt, Value: "5", Kind: listopt.KindInt}},
			wantSQL:  " AND player_number < $1",
			wantArgs: []any{5},
		},
		{
			name: "multiple filters increment placeholders",
			filters: []listopt.Filter{
				{Column: "min_players", Operator: listopt.Ge, Value: "3", Kind: listopt.KindInt},
				{Column: "name", Operator: listopt.Like, Value: "cat"},
			},
			wantSQL:  " AND min_players >= $1 AND name ILIKE $2",
			wantArgs: []any{3, "%cat%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotArgs := filterConditions(tt.filters, 1)
			require.Equal(t, tt.wantSQL, gotSQL)
			require.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func TestBuildPaginatedQuery(t *testing.T) {
	params := listopt.Apply(
		listopt.WithFilter(listopt.Filter{Column: "name", Operator: listopt.Like, Value: "cat"}),
		listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"}),
		listopt.WithPagination(listopt.Pagination{Page: 2, PageSize: 20}),
	)

	q, args := buildPaginatedQuery("SELECT * FROM boardgames WHERE deleted_at IS NULL", params)

	require.Contains(t, q, " AND name ILIKE $1")
	require.Contains(t, q, " ORDER BY name asc")
	require.Contains(t, q, " LIMIT $2 OFFSET $3")
	require.Equal(t, []any{"%cat%", 20, 20}, args) // limit=20, offset=(2-1)*20=20
}

func TestBuildPaginatedQuery_DefaultOrderByIDWhenNoSort(t *testing.T) {
	params := listopt.Apply(listopt.WithPagination(listopt.Pagination{Page: 1, PageSize: 10}))
	q, _ := buildPaginatedQuery("SELECT * FROM boardgames WHERE deleted_at IS NULL", params)

	require.Contains(t, q, " ORDER BY id ASC")
}

func TestBuildPaginatedQuery_IDTiebreakWithSort(t *testing.T) {
	params := listopt.Apply(listopt.WithSort(listopt.Sort{Column: "name", Order: "asc"}), listopt.WithPagination(listopt.Pagination{Page: 1, PageSize: 10}))
	q, _ := buildPaginatedQuery("SELECT * FROM boardgames WHERE deleted_at IS NULL", params)

	require.Contains(t, q, " ORDER BY name asc, id ASC")
}

func TestBuildPaginatedQuery_SortByIDHasNoDuplicateTiebreak(t *testing.T) {
	params := listopt.Apply(listopt.WithSort(listopt.Sort{Column: "id", Order: "desc"}), listopt.WithPagination(listopt.Pagination{Page: 1, PageSize: 10}))
	q, _ := buildPaginatedQuery("SELECT * FROM boardgames WHERE deleted_at IS NULL", params)

	require.Contains(t, q, " ORDER BY id desc")
	require.NotContains(t, q, "id desc, id ASC")
}

func TestBuildCountQuery(t *testing.T) {
	params := listopt.Apply(listopt.WithFilter(listopt.Filter{Column: "name", Operator: listopt.Eq, Value: "Catan"}))
	q, args := buildCountQuery("SELECT COUNT(*) FROM boardgames WHERE deleted_at IS NULL", params)

	require.Equal(t, "SELECT COUNT(*) FROM boardgames WHERE deleted_at IS NULL AND name = $1", q)
	require.Equal(t, []any{"Catan"}, args)
}
