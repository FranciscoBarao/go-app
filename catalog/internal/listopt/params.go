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

// Sort holds the sorting parameters.
type Sort struct {
	Column string
	Order  string
}

// Filter holds the filtering parameters.
type Filter struct {
	Column string
	Op     Op
	Value  string
}

// Params is the query contract passed from the service layer to the database layer.
type Params struct {
	Sort   Sort
	Filter Filter
}
