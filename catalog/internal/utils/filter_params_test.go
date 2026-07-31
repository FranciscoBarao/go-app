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

func TestGetFilterFields(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		op          string
		value       string
		col         string
		wantOp      listopt.Op
		wantNumeric bool
		err         bool
	}{
		{"empty op is like", "name", "", "catan", "name", listopt.OpLike, false, false},
		{"explicit like", "name", "like", "catan", "name", listopt.OpLike, false, false},
		{"like value with dots", "name", "like", "foo.bar", "name", listopt.OpLike, false, false},
		{"eq value with dots", "name", "eq", "foo.bar.baz", "name", listopt.OpEq, false, false},
		{"eq numeric-looking value on string col", "name", "eq", "123", "name", listopt.OpEq, false, false},
		{"eq on numeric col", "playernumber", "eq", "3", "player_number", listopt.OpEq, true, false},
		{"numeric decimal value", "playernumber", "ge", "3.5", "player_number", listopt.OpGe, true, false},
		{"empty field", "", "eq", "x", "", "", false, true},
		{"empty value", "name", "eq", "", "", "", false, true},
		{"invalid op", "name", "bogus", "x", "", "", false, true},
		{"unknown field", "unknown", "eq", "x", "", "", false, true},
		{"like on non-string", "playernumber", "like", "3", "", "", false, true},
		{"numeric op on string", "name", "lt", "5", "", "", false, true},
		{"field with db:-", "tags", "eq", "x", "", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, op, val, numeric, err := GetFilterFields(testModel{}, tt.field, tt.op, tt.value)
			if (err != nil) != tt.err {
				t.Fatalf("err = %v, wantErr = %v", err, tt.err)
			}
			if !tt.err {
				if col != tt.col || op != tt.wantOp || val != tt.value || numeric != tt.wantNumeric {
					t.Errorf("got (%q, %q, %q, %t), want (%q, %q, %q, %t)", col, op, val, numeric, tt.col, tt.wantOp, tt.value, tt.wantNumeric)
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
