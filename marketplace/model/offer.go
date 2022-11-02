package model

import "time"

type Offer struct {
	// Automatic
	Uuid    string    `json:"uuid,omitempty" db:"uuid" valid:"-"`
	AddedAt time.Time `json:"-" db:"added_at" valid:"-"`

	// User information
	Username string `json:"username,omitempty" db:"username"`

	// Catalog information
	Type string `json:"type" db:"type" valid:"required, alphanum, maxstringlength(100)"`

	// Offer information
	Name  string  `json:"name" db:"name" valid:"required, alphanum, maxstringlength(100)"`
	Price float64 `json:"price" db:"price" valid:"required, float, range(0|1000)"`
}

type OfferUpdate struct {
	Name  string  `json:"name" valid:"required, alphanum, maxstringlength(100)"`
	Price float64 `json:"price" valid:"required, float, range(0|1000)"`
}

// Update
func (off *Offer) UpdateOffer(offer *OfferUpdate) {
	off.Name = offer.GetName()
	off.Price = offer.GetPrice()

}

// Getters for Offer Update
func (off *OfferUpdate) GetName() string {
	return off.Name
}

func (off *OfferUpdate) GetPrice() float64 {
	return off.Price
}

// Getters for Offer
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

func (off *Offer) GetUsername() string {
	return off.Username
}

// Set
func (off *Offer) SetId(uuid string) {
	off.Uuid = uuid
}

func (off *Offer) SetUsername(username string) {
	off.Username = username
}

// Get Schema
func GetOfferSchema() string {
	var schema = `
	CREATE TABLE IF NOT EXISTS Offer (
			uuid uuid DEFAULT gen_random_uuid (),
			username text NOT NULL,
			type text,
			name text,
			price float,
			added_at timestamp DEFAULT now(),
			PRIMARY KEY (uuid)
		);`

	return schema
}
