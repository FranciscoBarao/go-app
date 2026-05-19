package utils

import (
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
)

type testModel struct {
	Name         string `db:"name"`
	PlayerNumber int    `db:"player_number"`
	Tags         []int  `db:"-"`
}

func TestGetFilter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		col   string
		op    listopt.Op
		val   string
		err   bool
	}{
		{"empty input", "", "", "", "", false},
		{"partial match", "name.catan", "name", listopt.OpLike, "catan", false},
		{"exact equality", "name.eq.Catan", "name", listopt.OpEq, "Catan", false},
		{"numeric lt", "playernumber.lt.5", "player_number", listopt.OpLt, "5", false},
		{"numeric le", "playernumber.le.10", "player_number", listopt.OpLe, "10", false},
		{"numeric gt", "playernumber.gt.2", "player_number", listopt.OpGt, "2", false},
		{"numeric ge", "playernumber.ge.3", "player_number", listopt.OpGe, "3", false},
		{"too many parts", "a.b.c.d", "", "", "", true},
		{"too few parts", "name", "", "", "", true},
		{"empty field", ".value", "", "", "", true},
		{"empty value", "name.", "", "", "", true},
		{"empty operator 3-part", "name..a", "", "", "", true},
		{"unknown field", "unknown.asc", "", "", "", true},
		{"invalid operator", "playernumber.xx.5", "", "", "", true},
		{"numeric op on string", "name.lt.5", "", "", "", true},
		{"non-numeric value", "playernumber.lt.abc", "", "", "", true},
		{"like on non-string", "playernumber.hello", "", "", "", true},
		{"field with db:-", "tags.test", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, op, val, err := GetFilter(testModel{}, tt.input)
			if (err != nil) != tt.err {
				t.Fatalf("err = %v, wantErr = %v", err, tt.err)
			}
			if !tt.err {
				if col != tt.col || op != tt.op || val != tt.val {
					t.Errorf("got (%q, %q, %q), want (%q, %q, %q)", col, op, val, tt.col, tt.op, tt.val)
				}
			}
		})
	}
}

func TestResolveField(t *testing.T) {
	tests := []struct {
		name  string
		field string
		col   string
		typ   string
		err   bool
	}{
		{"valid string field", "name", "name", "string", false},
		{"valid int field", "playernumber", "player_number", "int", false},
		{"case insensitive", "PlayerNumber", "player_number", "int", false},
		{"field with db:-", "tags", "", "", true},
		{"unknown field", "nonexistent", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := resolveField(testModel{}, tt.field)
			if (err != nil) != tt.err {
				t.Fatalf("err = %v, wantErr = %v", err, tt.err)
			}
			if !tt.err {
				if info.Column != tt.col || info.Type != tt.typ {
					t.Errorf("got (%q, %q), want (%q, %q)", info.Column, info.Type, tt.col, tt.typ)
				}
			}
		})
	}
}
