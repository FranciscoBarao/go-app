package listopt

// Op represents a validated filter operator.
type Op string

// Filter operators.
const (
	OpLike Op = "like"
	OpEq   Op = "eq"
	OpLt   Op = "lt"
	OpLe   Op = "le"
	OpGt   Op = "gt"
	OpGe   Op = "ge"
)

// Pagination defaults and bounds.
const (
	// DefaultPage is the page returned when none is requested.
	DefaultPage = 1
	// DefaultPageSize is the number of items per page when none is requested.
	DefaultPageSize = 10
	// MaxPageSize is the largest page size a client may request.
	MaxPageSize = 100
)

// Sort holds the sorting parameters.
type Sort struct {
	Column string
	Order  string
}

// Filter holds a single filtering clause.
type Filter struct {
	Column string
	Op     Op
	Value  string
	// Numeric indicates the target column is a numeric type, so the database
	// layer should bind Value as a number rather than a string.
	Numeric bool
}

// Pagination holds validated pagination parameters.
type Pagination struct {
	Page     int
	PageSize int
}

// Limit returns the SQL LIMIT (page size) for the pagination window.
func (p Pagination) Limit() int {
	return p.PageSize
}

// Offset returns the SQL OFFSET for the pagination window.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Params is the query contract passed from the service layer to the database layer.
type Params struct {
	Sort       Sort
	Filters    []Filter
	Pagination Pagination
}
