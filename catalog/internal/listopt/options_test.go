package listopt

import "testing"

func TestWithSort(t *testing.T) {
	p := Apply(WithSort("name", "asc"))
	if p.Sort.Column != "name" || p.Sort.Order != "asc" {
		t.Fatalf("expected name/asc, got %s/%s", p.Sort.Column, p.Sort.Order)
	}
}

func TestWithFilter(t *testing.T) {
	p := Apply(WithFilter("player_number", OpLt, "5"))
	if p.Filter.Column != "player_number" || p.Filter.Op != OpLt || p.Filter.Value != "5" {
		t.Fatalf("expected player_number/lt/5, got %s/%s/%s", p.Filter.Column, p.Filter.Op, p.Filter.Value)
	}
}

func TestApplySortAndFilter(t *testing.T) {
	p := Apply(WithSort("name", "asc"), WithFilter("player_number", OpGt, "2"))
	if p.Sort.Column != "name" || p.Sort.Order != "asc" {
		t.Fatalf("sort: expected name/asc, got %s/%s", p.Sort.Column, p.Sort.Order)
	}
	if p.Filter.Column != "player_number" || p.Filter.Op != OpGt || p.Filter.Value != "2" {
		t.Fatalf("filter: expected player_number/gt/2, got %s/%s/%s", p.Filter.Column, p.Filter.Op, p.Filter.Value)
	}
}

func TestApplyNoOptions(t *testing.T) {
	p := Apply()
	if p.Sort.Column != "" || p.Sort.Order != "" || p.Filter.Column != "" {
		t.Fatal("expected zero-value Params")
	}
}
