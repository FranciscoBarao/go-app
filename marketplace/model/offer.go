package model

type Offer struct {
	Name string `db:"name"`
}

// Constructors
func NewOffer(name string) Offer {
	return Offer{
		Name: name,
	}
}

// Update
func (off *Offer) UpdateOffer(offer Offer) {
	off.Name = offer.GetName()
}

// Getters
func (off Offer) GetName() string {
	return off.Name
}

// Get Schema
func GetOfferSchema() string {
	var schema = `
		CREATE TABLE Offer (
			name text,
			added_at timestamp default now()
		);`

	return schema
}
