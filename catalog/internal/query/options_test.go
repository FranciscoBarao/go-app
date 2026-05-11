package query

import "testing"

func TestWithSort(t *testing.T) {
	o := &Options{}
	WithSort("name", "asc")(o)
	if o.SortColumn != "name" || o.SortOrder != "asc" {
		t.Fatalf("expected name/asc, got %s/%s", o.SortColumn, o.SortOrder)
	}
}

func TestNoOptions(t *testing.T) {
	o := &Options{}
	if o.SortColumn != "" || o.SortOrder != "" {
		t.Fatal("expected zero values")
	}
}

func TestApply(t *testing.T) {
	f := Apply(WithSort("created_at", "desc"))
	if f.SortColumn != "created_at" || f.SortOrder != "desc" {
		t.Fatalf("expected created_at/desc, got %s/%s", f.SortColumn, f.SortOrder)
	}
}

func TestApplyNoOptions(t *testing.T) {
	f := Apply()
	if f.SortColumn != "" || f.SortOrder != "" {
		t.Fatal("expected zero-value Filter")
	}
}
