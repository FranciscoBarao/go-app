package mechanism

// CreateMechanismRequest is the request body for POST /mechanism.
type CreateMechanismRequest struct {
	Name  string `json:"name" valid:"required,maxstringlength(100)"`
	BggID *int   `json:"bgg_id,omitempty"`
}
