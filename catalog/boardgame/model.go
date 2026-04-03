package boardgame

import (
	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/tag"
	"gorm.io/gorm"
)

// Boardgame represents a board game entity with its associations.
type Boardgame struct {
	gorm.Model   `swaggerignore:"true"`
	Name         string                `json:"name" valid:"alphanum, maxstringlength(100)"`
	Publisher    string                `json:"publisher" valid:"alphanum, maxstringlength(100)"`
	PlayerNumber int                   `json:"playerNumber" valid:"int, range(1|16)"`
	Tags         []tag.Tag             `gorm:"many2many:boardgame_tags;" json:"tags,omitempty"`
	Categories   []category.Category   `gorm:"many2many:boardgame_categories;" json:"categories,omitempty"`
	Mechanisms   []mechanism.Mechanism `gorm:"many2many:boardgame_mechanisms;" json:"mechanisms,omitempty"`
	Ratings      []Rating              `gorm:"many2many:boardgame_ratings;" json:"ratings,omitempty"`
	Expansions   []Boardgame           `gorm:"foreignkey:BoardgameID" swaggerignore:"true" json:"expansions,omitempty"`
	BoardgameID  *uint                 `swaggerignore:"true" json:"boardgame_id,omitempty"`
}

// UpdateBoardgame applies changes from the given boardgame to this one.
func (bg *Boardgame) UpdateBoardgame(boardgame *Boardgame) {
	bg.Name = boardgame.GetName()
	bg.Publisher = boardgame.GetPublisher()
	bg.PlayerNumber = boardgame.GetPlayerNumber()
	bg.Tags = boardgame.GetTags()
	bg.Categories = boardgame.GetCategories()
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

// GetID returns the boardgame's primary key.
func (bg Boardgame) GetID() *uint {
	return &bg.ID
}

// GetName returns the boardgame's name.
func (bg Boardgame) GetName() string {
	return bg.Name
}

// GetPublisher returns the boardgame's publisher.
func (bg Boardgame) GetPublisher() string {
	return bg.Publisher
}

// GetPlayerNumber returns the boardgame's player count.
func (bg Boardgame) GetPlayerNumber() int {
	return bg.PlayerNumber
}

// GetTags returns the boardgame's tags.
func (bg Boardgame) GetTags() []tag.Tag {
	return bg.Tags
}

// GetCategories returns the boardgame's categories.
func (bg Boardgame) GetCategories() []category.Category {
	return bg.Categories
}

// GetMechanisms returns the boardgame's mechanisms.
func (bg Boardgame) GetMechanisms() []mechanism.Mechanism {
	return bg.Mechanisms
}

// GetExpansions returns the boardgame's expansions.
func (bg Boardgame) GetExpansions() []Boardgame {
	return bg.Expansions
}

// GetBoardgameID returns the parent boardgame ID if this is an expansion.
func (bg Boardgame) GetBoardgameID() *uint {
	return bg.BoardgameID
}

// SetBoardgameID sets the parent boardgame ID for an expansion.
func (bg *Boardgame) SetBoardgameID(id *uint) {
	bg.BoardgameID = id
}
