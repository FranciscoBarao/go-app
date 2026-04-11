package boardgame

import (
	"time"

	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/tag"
)

// Boardgame represents a board game entity with its associations.
type Boardgame struct {
	ID           uint                  `json:"id"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	Name         string                `json:"name" valid:"alphanum, maxstringlength(100)"`
	Publisher    string                `json:"publisher" valid:"alphanum, maxstringlength(100)"`
	PlayerNumber int                   `json:"playerNumber" valid:"int, range(1|16)"`
	Tags         []tag.Tag             `json:"tags,omitempty"`
	Categories   []category.Category   `json:"categories,omitempty"`
	Mechanisms   []mechanism.Mechanism `json:"mechanisms,omitempty"`
	Ratings      []Rating              `json:"ratings,omitempty"`
	Expansions   []Boardgame           `json:"expansions,omitempty"`
	BoardgameID  *uint                 `json:"boardgame_id,omitempty"`
}

// UpdateBoardgame applies changes from the given boardgame to this one.
func (bg *Boardgame) UpdateBoardgame(input *Boardgame) {
	bg.Name = input.Name
	bg.Publisher = input.Publisher
	bg.PlayerNumber = input.PlayerNumber
	bg.Tags = input.Tags
	bg.Categories = input.Categories
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
