package database

import (
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
)

func TestFilterClause(t *testing.T) {
	tests := []struct {
		name     string
		params   listopt.Params
		wantSQL  string
		wantArg  any
	}{
		{
			name:    "empty filter returns nothing",
			params:  listopt.Params{},
			wantSQL: "",
			wantArg: nil,
		},
		{
			name:    "like op produces ILIKE clause",
			params:  listopt.Params{Filter: listopt.Filter{Column: "name", Op: listopt.OpLike, Value: "catan"}},
			wantSQL: " WHERE name ILIKE $1",
			wantArg: "%catan%",
		},
		{
			name:    "eq op produces equality clause",
			params:  listopt.Params{Filter: listopt.Filter{Column: "name", Op: listopt.OpEq, Value: "Catan"}},
			wantSQL: " WHERE name = $1",
			wantArg: "Catan",
		},
		{
			name:    "lt op with numeric value",
			params:  listopt.Params{Filter: listopt.Filter{Column: "player_number", Op: listopt.OpLt, Value: "5"}},
			wantSQL: " WHERE player_number < $1",
			wantArg: 5,
		},
		{
			name:    "ge op with float value",
			params:  listopt.Params{Filter: listopt.Filter{Column: "rating", Op: listopt.OpGe, Value: "3.5"}},
			wantSQL: " WHERE rating >= $1",
			wantArg: 3.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotArg := filterClause(tt.params)
			if gotSQL != tt.wantSQL {
				t.Errorf("sql = %q, want %q", gotSQL, tt.wantSQL)
			}
			if gotArg != tt.wantArg {
				t.Errorf("arg = %v (%T), want %v (%T)", gotArg, gotArg, tt.wantArg, tt.wantArg)
			}
		})
	}
}
