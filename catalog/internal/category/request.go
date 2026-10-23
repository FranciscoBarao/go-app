package category

// CreateCategoryRequest is the request body for POST /category.
type CreateCategoryRequest struct {
	Name  string `json:"name" valid:"required,maxstringlength(100)"`
	BggID *int   `json:"bgg_id,omitempty"`
}
