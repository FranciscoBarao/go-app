package model

import "time"

type Offer struct {
	// Automatic
	Uuid    string    `json:"uuid,omitempty" db:"uuid" valid:"-"`
	AddedAt time.Time `json:"-" db:"added_at" valid:"-"`

	// Catalog information
	Type string `json:"type" db:"type" valid:"required, alphanum, maxstringlength(100)"`

	// Offer information
	Name  string  `json:"name" db:"name" valid:"required, alphanum, maxstringlength(100)"`
	Price float64 `json:"price" db:"price" valid:"required, float, range(0|1000)"`
}

// Constructors
func NewOffer(name string, price float64) Offer {
	return Offer{
		Name:  name,
		Price: price,
	}
}

// Update
func (off *Offer) UpdateOffer(offer *Offer) {
	off.Name = offer.GetName()
	off.Price = offer.GetPrice()

}

// Getters
func (off *Offer) GetName() string {
	return off.Name
}

func (off *Offer) GetType() string {
	return off.Type
}

func (off *Offer) GetPrice() float64 {
	return off.Price
}

func (off *Offer) GetId() string {
	return off.Uuid
}

// Set
func (off *Offer) SetId(uuid string) {
	off.Uuid = uuid
}

// Get Schema
func GetOfferSchema() string {
	var schema = `
	CREATE TABLE IF NOT EXISTS Offer (
			uuid uuid DEFAULT gen_random_uuid (),
			type text,
			name text,
			price float,
			added_at timestamp DEFAULT now(),
			PRIMARY KEY (uuid)
		);`

	return schema
}
