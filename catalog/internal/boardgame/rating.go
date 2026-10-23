package boardgame

import "time"

// Rating represents a user's rating for a boardgame.
type Rating struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username,omitempty" valid:"required"`
	Value     int       `json:"value" valid:"required, int, range(0|10)"`
}
