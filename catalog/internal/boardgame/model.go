package boardgame

import (
	"time"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/tag"
)

// Boardgame represents a board game entity with its associations.
type Boardgame struct {
	ID           uint                  `json:"id" db:"id"`
	CreatedAt    time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at" db:"updated_at"`
	Name         string                `json:"name" db:"name" valid:"alphanum, maxstringlength(100)"`
	Publisher    string                `json:"publisher" db:"publisher" valid:"alphanum, maxstringlength(100)"`
	PlayerNumber int                   `json:"playerNumber" db:"player_number" valid:"int, range(1|16)"`
	Tags         []tag.Tag             `json:"tags,omitempty" db:"-"`
	Categories   []category.Category   `json:"categories,omitempty" db:"-"`
	Mechanisms   []mechanism.Mechanism `json:"mechanisms,omitempty" db:"-"`
	Ratings      []Rating              `json:"ratings,omitempty" db:"-"`
	Expansions   []Boardgame           `json:"expansions,omitempty" db:"-"`
	BoardgameID  *uint                 `json:"boardgame_id,omitempty" db:"-"`
}

// UpdateBoardgame applies changes from the given boardgame to this one.
func (bg *Boardgame) UpdateBoardgame(input *Boardgame) {
	bg.Name = input.Name
	bg.Publisher = input.Publisher
	bg.PlayerNumber = input.PlayerNumber
	bg.Tags = input.Tags
	bg.Categories = input.Categories
	bg.Mechanisms = input.Mechanisms
}

// HasTags returns true if the boardgame has associated tags.
func (bg Boardgame) HasTags() bool {
	return len(bg.Tags) > 0
}

// HasCategories returns true if the boardgame has associated categories.
func (bg Boardgame) HasCategories() bool {
	return len(bg.Categories) > 0
}

// HasMechanisms returns true if the boardgame has associated mechanisms.
func (bg Boardgame) HasMechanisms() bool {
	return len(bg.Mechanisms) > 0
}

// HasExpansions returns true if the boardgame has expansions.
func (bg Boardgame) HasExpansions() bool {
	return len(bg.Expansions) > 0
}

// IsExpansion returns true if this boardgame is an expansion of another.
func (bg Boardgame) IsExpansion() bool {
	return bg.BoardgameID != nil
}
