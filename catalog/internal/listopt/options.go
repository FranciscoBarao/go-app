package listopt

// Option is a functional option for configuring list query parameters.
type Option func(*Params)

// WithSort sets the sort column and order.
func WithSort(s Sort) Option {
	return func(params *Params) {
		params.Sort = s
	}
}

// WithFilter appends a filter clause. Multiple filters are combined with AND at
// the database layer.
func WithFilter(f Filter) Option {
	return func(params *Params) {
		params.Filters = append(params.Filters, f)
	}
}

// WithPagination sets the pagination window, clamping page/pageSize into bounds.
func WithPagination(p Pagination) Option {
	return func(params *Params) {
		params.Pagination = p.Normalize()
	}
}

// Apply processes options and returns Params ready for the database layer.
// Pagination defaults to DefaultPage / DefaultPageSize
func Apply(opts ...Option) Params {
	p := Params{
		Pagination: NewPagination(DefaultPage, DefaultPageSize),
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}
