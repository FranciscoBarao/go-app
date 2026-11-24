package listopt

// Option is a functional option for configuring list query parameters.
type Option func(*Query)

// WithSort sets the sort column and order.
func WithSort(sort Sort) Option {
	return func(query *Query) {
		query.Sort = sort
	}
}

// WithFilter appends a filter clause. Multiple filters are combined with AND at
// the database layer.
func WithFilter(f Filter) Option {
	return func(query *Query) {
		query.Filters = append(query.Filters, f)
	}
}

// WithPagination sets the pagination window, clamping page/pageSize into bounds.
func WithPagination(p Pagination) Option {
	return func(query *Query) {
		query.Pagination = p.Normalize()
	}
}

// NewQuery processes options and returns Query ready for the database layer.
// Pagination defaults to DefaultPage / DefaultPageSize.
func NewQuery(opts ...Option) Query {
	query := Query{
		Pagination: NewPagination(DefaultPage, DefaultPageSize),
	}
	for _, opt := range opts {
		opt(&query)
	}
	return query
}
