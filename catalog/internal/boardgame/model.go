package boardgame

import (
	"time"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
)

// Boardgame represents a board game entity with its associations.
type Boardgame struct {
	ID            uint                       `json:"id" db:"id"`
	Slug          string                     `json:"slug" db:"slug"`
	CreatedAt     time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time                 `json:"deleted_at,omitempty" db:"deleted_at"`
	Name          string                     `json:"name" db:"name" valid:"required,maxstringlength(120)"`
	Description   string                     `json:"description" db:"description"`
	YearPublished int                        `json:"year_published" db:"year_published"`
	MinPlayers    int                        `json:"min_players" db:"min_players" valid:"required,range(1|16)"`
	MaxPlayers    int                        `json:"max_players" db:"max_players" valid:"required,range(1|16)"`
	MinPlayTime   int                        `json:"min_play_time" db:"min_play_time"`
	MaxPlayTime   int                        `json:"max_play_time" db:"max_play_time"`
	MinAge        *int                       `json:"min_age,omitempty" db:"min_age"`
	BggID         *int                       `json:"bgg_id,omitempty" db:"bgg_id"`
	Categories    []category.Category        `json:"categories,omitempty" db:"-"`
	Mechanisms    []mechanism.Mechanism      `json:"mechanisms,omitempty" db:"-"`
	Contributions []contributor.Contribution `json:"contributions,omitempty" db:"-"`
	Ratings       []Rating                   `json:"ratings,omitempty" db:"-"`
	Expansions    []Boardgame                `json:"expansions,omitempty" db:"-"`
	BoardgameID   *uint                      `json:"boardgame_id,omitempty" db:"boardgame_id"`
}

// HasCategories returns true if the boardgame has associated categories.
func (bg Boardgame) HasCategories() bool {
	return len(bg.Categories) > 0
}

// HasMechanisms returns true if the boardgame has associated mechanisms.
func (bg Boardgame) HasMechanisms() bool {
	return len(bg.Mechanisms) > 0
}

// HasContributions returns true if the boardgame has contributor credits.
func (bg Boardgame) HasContributions() bool {
	return len(bg.Contributions) > 0
}

// HasExpansions returns true if the boardgame has expansions.
func (bg Boardgame) HasExpansions() bool {
	return len(bg.Expansions) > 0
}

// IsExpansion returns true if this boardgame is an expansion of another.
func (bg Boardgame) IsExpansion() bool {
	return bg.BoardgameID != nil
}

// IsDeleted returns true when the boardgame has been soft-deleted.
func (bg Boardgame) IsDeleted() bool {
	return bg.DeletedAt != nil
}
