package listopt

import "testing"

func testSchema() Schema {
	return Schema{
		"name":        StringField("name", Sortable, Filterable),
		"count":       IntField("count", Sortable, Filterable),
		"notes":       StringField("notes", NotSortable, Filterable),
		"created_at":  SortColumn("created_at"),
		"internal_id": {Column: "internal_id", Kind: KindInt, Sortable: NotSortable, Filterable: NotFilterable},
	}
}

func TestSchemaSort(t *testing.T) {
	s := testSchema()
	tests := []struct {
		name      string
		field     string
		order     string
		col       string
		wantOrder string
		err       bool
	}{
		{"empty", "", "", "", "", false},
		{"asc", "name", "asc", "name", "asc", false},
		{"desc", "count", "desc", "count", "desc", false},
		{"sort only", "created_at", "asc", "created_at", "asc", false},
		{"empty field", "", "asc", "", "", true},
		{"empty order", "name", "", "", "", true},
		{"bad order", "name", "sideways", "", "", true},
		{"unknown", "missing", "asc", "", "", true},
		{"filter only", "notes", "asc", "", "", true},
		{"omitted", "internal_id", "asc", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Sort(tt.field, tt.order)
			if (err != nil) != tt.err {
				t.Fatalf("err = %v, wantErr = %v", err, tt.err)
			}
			if !tt.err && (got.Column != tt.col || got.Order != tt.wantOrder) {
				t.Errorf("got (%q, %q), want (%q, %q)", got.Column, got.Order, tt.col, tt.wantOrder)
			}
		})
	}
}

func TestSchemaFilter(t *testing.T) {
	s := testSchema()
	tests := []struct {
		name         string
		field        string
		op           string
		value        string
		col          string
		wantOperator Operator
		wantKind     FieldKind
		err          bool
	}{
		{"empty op is like", "name", "", "catan", "name", Like, KindString, false},
		{"explicit like", "name", "like", "foo.bar", "name", Like, KindString, false},
		{"eq string numeric-looking", "name", "eq", "123", "name", Eq, KindString, false},
		{"eq int", "count", "eq", "3", "count", Eq, KindInt, false},
		{"ge decimal", "count", "ge", "3.5", "count", Ge, KindInt, false},
		{"filter only", "notes", "like", "x", "notes", Like, KindString, false},
		{"empty field", "", "eq", "x", "", "", 0, true},
		{"empty value", "name", "eq", "", "", "", 0, true},
		{"invalid op", "name", "bogus", "x", "", "", 0, true},
		{"unknown", "missing", "eq", "x", "", "", 0, true},
		{"like on int", "count", "like", "3", "", "", 0, true},
		{"numeric op on string", "name", "lt", "5", "", "", 0, true},
		{"non-numeric value", "count", "lt", "abc", "", "", 0, true},
		{"eq int non-numeric", "count", "eq", "abc", "", "", 0, true},
		{"sort only", "created_at", "eq", "x", "", "", 0, true},
		{"omitted", "internal_id", "eq", "1", "", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Filter(tt.field, tt.op, tt.value)
			if (err != nil) != tt.err {
				t.Fatalf("err = %v, wantErr = %v", err, tt.err)
			}
			if !tt.err {
				if got.Column != tt.col || got.Operator != tt.wantOperator || got.Value != tt.value || got.Kind != tt.wantKind {
					t.Errorf("got (%q, %q, %q, %v), want (%q, %q, %q, %v)",
						got.Column, got.Operator, got.Value, got.Kind, tt.col, tt.wantOperator, tt.value, tt.wantKind)
				}
			}
		})
	}
}
