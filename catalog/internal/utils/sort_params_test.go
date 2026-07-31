package utils

import "testing"

func TestGetSort(t *testing.T) {
	tests := []struct {
		name  string
		input string
		col   string
		order string
		err   bool
	}{
		{"empty input", "", "", "", false},
		{"asc", "name.asc", "name", "asc", false},
		{"desc", "name.desc", "name", "desc", false},
		{"case insensitive field", "PlayerNumber.asc", "player_number", "asc", false},
		{"too few parts", "name", "", "", true},
		{"too many parts", "name.asc.extra", "", "", true},
		{"empty field", ".asc", "", "", true},
		{"empty order", "name.", "", "", true},
		{"invalid order", "name.random", "", "", true},
		{"unknown field", "unknown.asc", "", "", true},
		{"field with db:-", "tags.asc", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, order, err := GetSort(testModel{}, tt.input)
			if (err != nil) != tt.err {
				t.Fatalf("err = %v, wantErr = %v", err, tt.err)
			}
			if !tt.err {
				if col != tt.col || order != tt.order {
					t.Errorf("got (%q, %q), want (%q, %q)", col, order, tt.col, tt.order)
				}
			}
		})
	}
}
