package boardgame

import "gorm.io/gorm"

// Rating struct - co-located since Rating is boardgame-only
type Rating struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Username string `json:"username,omitempty" db:"username" gorm:"unique"`
	Value    int    `json:"value" db:"value" valid:"required, int, range(0|10)"`
}

func (rating *Rating) SetUsername(username string) {
	rating.Username = username
}
