package contributor

import "time"

// Role identifies the type of credit on a boardgame.
type Role string

// Supported contributor roles on a boardgame credit.
const (
	RoleDesigner        Role = "designer"
	RoleArtist          Role = "artist"
	RolePublisher       Role = "publisher"
	RoleDeveloper       Role = "developer"
	RoleGraphicDesigner Role = "graphic_designer"
)

// Contributor is a boardgame industry credit (designer, publisher, etc).
type Contributor struct {
	ID        uint       `json:"id" db:"id"`
	Slug      string     `json:"slug,omitempty" db:"slug" valid:"optional,maxstringlength(100)"`
	Name      string     `json:"name" db:"name" valid:"required,maxstringlength(100)"`
	Bio       *string    `json:"bio,omitempty" db:"bio"`
	BggID     *int       `json:"bgg_id,omitempty" db:"bgg_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Contribution links a contributor to a boardgame with a role.
type Contribution struct {
	Contributor Contributor `json:"contributor"`
	Role        Role        `json:"role" valid:"required"`
	CreditOrder *int        `json:"credit_order,omitempty"`
}

// ContributionInput is used when creating/updating boardgame credits.
type ContributionInput struct {
	Slug        string `json:"slug" valid:"required,maxstringlength(100)"`
	Role        Role   `json:"role" valid:"required"`
	CreditOrder *int   `json:"credit_order,omitempty"`
}

// Valid reports whether role is a supported credit type.
func (r Role) Valid() bool {
	switch r {
	case RoleDesigner, RoleArtist, RolePublisher, RoleDeveloper, RoleGraphicDesigner:
		return true
	default:
		return false
	}
}
