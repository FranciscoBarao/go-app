package listopt

// Option is a functional option for configuring list query parameters.
type Option func(*Params)

// WithSort sets the sort column and order.
func WithSort(column, order string) Option {
	return func(p *Params) {
		p.Sort = Sort{Column: column, Order: order}
	}
}

// WithFilter appends a filter clause (column, operator, value). The numeric flag
// marks the column as a numeric type so the database layer binds the value as a
// number rather than a string. Multiple filters are combined with AND at the
// database layer.
func WithFilter(column string, op Op, value string, numeric bool) Option {
	return func(p *Params) {
		p.Filters = append(p.Filters, Filter{Column: column, Op: op, Value: value, Numeric: numeric})
	}
}

// WithPagination sets the pagination window (page, pageSize).
func WithPagination(page, pageSize int) Option {
	return func(p *Params) {
		p.Pagination = Pagination{Page: page, PageSize: pageSize}
	}
}

// Apply processes options, normalizes pagination, and returns Params ready for
// the database layer.
func Apply(opts ...Option) Params {
	p := Params{}
	for _, opt := range opts {
		opt(&p)
	}
	p.Pagination = normalizePagination(p.Pagination)
	return p
}

// normalizePagination clamps pagination into valid bounds, applying defaults
// for unset or out-of-range values.
func normalizePagination(p Pagination) Pagination {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	switch {
	case p.PageSize < 1:
		p.PageSize = DefaultPageSize
	case p.PageSize > MaxPageSize:
		p.PageSize = MaxPageSize
	}
	return p
}
