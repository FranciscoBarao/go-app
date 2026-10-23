package offer

import "time"

// Offer represents a marketplace offer entity.
type Offer struct {
	UUID      string    `json:"uuid,omitempty" db:"uuid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Username  string    `json:"username,omitempty" db:"username"`
	Type      string    `json:"type" db:"type"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
}
