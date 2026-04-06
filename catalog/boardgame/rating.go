package boardgame

import "gorm.io/gorm"

// Rating represents a user's rating for a boardgame.
type Rating struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Username string `json:"username,omitempty" db:"username" gorm:"unique"`
	Value    int    `json:"value" db:"value" valid:"required, int, range(0|10)"`
}
