package query

// Filter is the query contract passed from the service layer to the database layer.
type Filter struct {
	SortColumn string
	SortOrder  string
}
