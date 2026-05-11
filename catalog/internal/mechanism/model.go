package mechanism

import "time"

// Mechanism represents a board game mechanism.
type Mechanism struct {
	Name      string    `json:"name" valid:"alphanum, maxstringlength(30)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewMechanism creates a new Mechanism with the given name.
func NewMechanism(name string) *Mechanism {
	return &Mechanism{
		Name: name,
	}
}
