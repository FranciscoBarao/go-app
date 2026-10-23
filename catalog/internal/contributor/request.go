package contributor

// CreateContributorRequest is the request body for POST /contributors.
type CreateContributorRequest struct {
	Name  string  `json:"name" valid:"required,maxstringlength(100)"`
	Bio   *string `json:"bio,omitempty"`
	BggID *int    `json:"bgg_id,omitempty"`
}

// UpdateContributorRequest is the request body for PATCH /contributors/{slug}.
type UpdateContributorRequest struct {
	Name  *string `json:"name,omitempty" valid:"optional,maxstringlength(100)"`
	Bio   *string `json:"bio,omitempty"`
	BggID *int    `json:"bgg_id,omitempty"`
}

// ToContributor applies non-nil fields from the request onto an existing Contributor.
func (r *UpdateContributorRequest) ToContributor(c *Contributor) {
	if r.Name != nil {
		c.Name = *r.Name
	}
	if r.Bio != nil {
		c.Bio = r.Bio
	}
	if r.BggID != nil {
		c.BggID = r.BggID
	}
}
