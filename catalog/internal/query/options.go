package query

// Option is a functional option for configuring query parameters.
type Option func(*Options)

// Options holds the configured query parameters.
type Options struct {
	SortColumn string
	SortOrder  string
}

// WithSort sets the sort column and order.
func WithSort(column, order string) Option {
	return func(o *Options) {
		o.SortColumn = column
		o.SortOrder = order
	}
}

// Apply processes options and returns a Filter ready for the database layer.
func Apply(opts ...Option) Filter {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	return Filter{SortColumn: o.SortColumn, SortOrder: o.SortOrder}
}
