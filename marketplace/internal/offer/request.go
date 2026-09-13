package offer

// CreateOfferRequest is the request DTO for POST /offer.
type CreateOfferRequest struct {
	Type  string  `json:"type" valid:"required, alphanum, maxstringlength(100)"`
	Name  string  `json:"name" valid:"required, alphanum, maxstringlength(100)"`
	Price float64 `json:"price" valid:"required, float, range(0|1000)"`
}

// NewOffer creates an Offer domain model from a CreateOfferRequest.
func NewOffer(req *CreateOfferRequest) *Offer {
	return &Offer{
		Type:  req.Type,
		Name:  req.Name,
		Price: req.Price,
	}
}

// UpdateOfferRequest is the request DTO for PATCH /offer/{id}.
type UpdateOfferRequest struct {
	Name  string  `json:"name" valid:"required, alphanum, maxstringlength(100)"`
	Price float64 `json:"price" valid:"required, float, range(0|1000)"`
}

// ApplyTo applies the update fields to an existing Offer.
func (r *UpdateOfferRequest) ApplyTo(o *Offer) {
	o.Name = r.Name
	o.Price = r.Price
}
