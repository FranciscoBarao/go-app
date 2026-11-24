package mechanism

import "time"

// Mechanism represents a board game mechanism (how you play).
type Mechanism struct {
	ID        uint       `json:"id" db:"id"`
	Slug      string     `json:"slug,omitempty" db:"slug" valid:"optional,maxstringlength(100)"`
	Name      string     `json:"name" db:"name" valid:"required,maxstringlength(100)"`
	BggID     *int       `json:"bgg_id,omitempty" db:"bgg_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
