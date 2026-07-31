package listopt

import "testing"

func TestWithSort(t *testing.T) {
	p := Apply(WithSort("name", "asc"))
	if p.Sort.Column != "name" || p.Sort.Order != "asc" {
		t.Fatalf("expected name/asc, got %s/%s", p.Sort.Column, p.Sort.Order)
	}
}

func TestWithFilter(t *testing.T) {
	p := Apply(WithFilter("player_number", OpLt, "5", true))
	if len(p.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(p.Filters))
	}
	f := p.Filters[0]
	if f.Column != "player_number" || f.Op != OpLt || f.Value != "5" || !f.Numeric {
		t.Fatalf("expected player_number/lt/5/numeric, got %s/%s/%s/%t", f.Column, f.Op, f.Value, f.Numeric)
	}
}

func TestWithMultipleFilters(t *testing.T) {
	p := Apply(
		WithFilter("min_players", OpGe, "3", true),
		WithFilter("name", OpLike, "cat", false),
	)
	if len(p.Filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(p.Filters))
	}
	if p.Filters[0].Column != "min_players" || p.Filters[1].Column != "name" {
		t.Fatalf("filters out of order: %+v", p.Filters)
	}
}

func TestApplySortAndFilter(t *testing.T) {
	p := Apply(WithSort("name", "asc"), WithFilter("player_number", OpGt, "2", true))
	if p.Sort.Column != "name" || p.Sort.Order != "asc" {
		t.Fatalf("sort: expected name/asc, got %s/%s", p.Sort.Column, p.Sort.Order)
	}
	if len(p.Filters) != 1 || p.Filters[0].Op != OpGt {
		t.Fatalf("filter: expected one gt filter, got %+v", p.Filters)
	}
}

func TestApplyNoOptions(t *testing.T) {
	p := Apply()
	if p.Sort.Column != "" || p.Sort.Order != "" || len(p.Filters) != 0 {
		t.Fatal("expected zero-value sort/filters")
	}
	// Pagination should be normalized to defaults.
	if p.Pagination.Page != DefaultPage || p.Pagination.PageSize != DefaultPageSize {
		t.Fatalf("expected default pagination %d/%d, got %d/%d",
			DefaultPage, DefaultPageSize, p.Pagination.Page, p.Pagination.PageSize)
	}
}

func TestWithPagination(t *testing.T) {
	p := Apply(WithPagination(3, 25))
	if p.Pagination.Page != 3 || p.Pagination.PageSize != 25 {
		t.Fatalf("expected 3/25, got %d/%d", p.Pagination.Page, p.Pagination.PageSize)
	}
	if p.Pagination.Offset() != 50 || p.Pagination.Limit() != 25 {
		t.Fatalf("expected offset 50 limit 25, got offset %d limit %d",
			p.Pagination.Offset(), p.Pagination.Limit())
	}
}

func TestPaginationNormalization(t *testing.T) {
	cases := []struct {
		name             string
		page, size       int
		wantPage, wantSz int
	}{
		{"zero -> defaults", 0, 0, DefaultPage, DefaultPageSize},
		{"negative -> defaults", -5, -1, DefaultPage, DefaultPageSize},
		{"over max -> capped", 2, 500, 2, MaxPageSize},
		{"at max -> kept", 1, MaxPageSize, 1, MaxPageSize},
		{"valid -> kept", 4, 20, 4, 20},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := Apply(WithPagination(c.page, c.size))
			if p.Pagination.Page != c.wantPage || p.Pagination.PageSize != c.wantSz {
				t.Fatalf("got %d/%d, want %d/%d",
					p.Pagination.Page, p.Pagination.PageSize, c.wantPage, c.wantSz)
			}
		})
	}
}
