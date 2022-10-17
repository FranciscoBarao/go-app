package model

type Offer struct {
	Uuid  string  `db:"uuid"`
	Name  string  `db:"name"`
	Type  string  `db:"type"`
	Price float64 `db:"price"`
}

// Constructors
func NewOffer(name string, price float64) Offer {
	return Offer{
		Name:  name,
		Price: price,
	}
}

// Update
func (off *Offer) UpdateOffer(offer Offer) {
	off.Name = offer.GetName()
	off.Price = offer.GetPrice()

}

// Getters
func (off Offer) GetName() string {
	return off.Name
}

func (off Offer) GetType() string {
	return off.Type
}

func (off Offer) GetPrice() float64 {
	return off.Price
}

// Get Schema
func GetOfferSchema() string {
	var schema = `
		CREATE TABLE IF NOT EXISTS Offer (
			uuid uuid DEFAULT uuid_generate_v4 (),
			name text,
			price float,
			added_at timestamp DEFAULT now(),
			PRIMARY KEY (uuid)
		);`

	return schema
}
