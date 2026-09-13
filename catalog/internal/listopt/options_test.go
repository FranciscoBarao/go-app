package listopt

import "testing"

func TestWithSort(t *testing.T) {
	p := NewQuery(WithSort(Sort{Column: "name", Order: "asc"}))
	if p.Sort.Column != "name" || p.Sort.Order != "asc" {
		t.Fatalf("expected name/asc, got %s/%s", p.Sort.Column, p.Sort.Order)
	}
}

func TestWithFilter(t *testing.T) {
	p := NewQuery(WithFilter(Filter{Column: "player_number", Operator: Lt, Value: "5", ValueKind: KindInt}))
	if len(p.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(p.Filters))
	}
	f := p.Filters[0]
	if f.Column != "player_number" || f.Operator != Lt || f.Value != "5" || f.ValueKind != KindInt {
		t.Fatalf("expected player_number/lt/5/int, got %s/%s/%s/%v", f.Column, f.Operator, f.Value, f.ValueKind)
	}
}

func TestWithMultipleFilters(t *testing.T) {
	p := NewQuery(
		WithFilter(Filter{Column: "min_players", Operator: Ge, Value: "3", ValueKind: KindInt}),
		WithFilter(Filter{Column: "name", Operator: Like, Value: "cat"}),
	)
	if len(p.Filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(p.Filters))
	}
	if p.Filters[0].Column != "min_players" || p.Filters[1].Column != "name" {
		t.Fatalf("filters out of order: %+v", p.Filters)
	}
}

func TestNewQuerySortAndFilter(t *testing.T) {
	p := NewQuery(WithSort(Sort{Column: "name", Order: "asc"}), WithFilter(Filter{Column: "player_number", Operator: Gt, Value: "2", ValueKind: KindInt}))
	if p.Sort.Column != "name" || p.Sort.Order != "asc" {
		t.Fatalf("sort: expected name/asc, got %s/%s", p.Sort.Column, p.Sort.Order)
	}
	if len(p.Filters) != 1 || p.Filters[0].Operator != Gt {
		t.Fatalf("filter: expected one gt filter, got %+v", p.Filters)
	}
}

func TestNewQueryNoOptions(t *testing.T) {
	p := NewQuery()
	if p.Sort.Column != "" || p.Sort.Order != "" || len(p.Filters) != 0 {
		t.Fatal("expected zero-value sort/filters")
	}
	if p.Pagination.Page != DefaultPage || p.Pagination.PageSize != DefaultPageSize {
		t.Fatalf("expected default pagination %d/%d, got %d/%d",
			DefaultPage, DefaultPageSize, p.Pagination.Page, p.Pagination.PageSize)
	}
}

func TestWithPagination(t *testing.T) {
	p := NewQuery(WithPagination(Pagination{Page: 3, PageSize: 25}))
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
			p := NewQuery(WithPagination(Pagination{Page: c.page, PageSize: c.size}))
			if p.Pagination.Page != c.wantPage || p.Pagination.PageSize != c.wantSz {
				t.Fatalf("got %d/%d, want %d/%d",
					p.Pagination.Page, p.Pagination.PageSize, c.wantPage, c.wantSz)
			}
		})
	}
}
