package listopt

// Option is a functional option for configuring list query parameters.
type Option func(*Params)

// WithSort sets the sort column and order.
func WithSort(column, order string) Option {
	return func(p *Params) {
		p.Sort = Sort{Column: column, Order: order}
	}
}

// WithFilter sets the filter column, operator, and value.
func WithFilter(column string, op Op, value string) Option {
	return func(p *Params) {
		p.Filter = Filter{Column: column, Op: op, Value: value}
	}
}

// Apply processes options and returns Params ready for the database layer.
func Apply(opts ...Option) Params {
	p := Params{}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}
