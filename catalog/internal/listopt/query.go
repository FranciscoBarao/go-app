package listopt

// Operator is a validated filter comparison (like, eq, lt, le, gt, ge).
type Operator string

const (
	Like Operator = "like"
	Eq   Operator = "eq"
	Lt   Operator = "lt"
	Le   Operator = "le"
	Gt   Operator = "gt"
	Ge   Operator = "ge"
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

// Query is the query contract passed from the service layer to the database layer.
type Query struct {
	Sort       Sort
	Filters    []Filter
	Pagination Pagination
}

// Sort holds the sorting parameters.
type Sort struct {
	Column string
	Order  string
}

// Filter holds a single filtering clause.
type Filter struct {
	Column    string
	Operator  Operator
	Value     string
	ValueKind FieldKind
}

// Pagination holds validated pagination parameters.
type Pagination struct {
	Page     int
	PageSize int
}

// NewPagination builds a Pagination from request ints and normalizes it.
func NewPagination(page, pageSize int) Pagination {
	return Pagination{Page: page, PageSize: pageSize}.Normalize()
}

// Limit returns the SQL LIMIT (page size) for the pagination window.
func (p Pagination) Limit() int {
	return p.PageSize
}

// Offset returns the SQL OFFSET for the pagination window.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Normalize clamps page/pageSize into valid bounds (defaults and max size).
func (p Pagination) Normalize() Pagination {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
	return p
}
