package model

import (
	"gorm.io/gorm"
)

type Boardgame struct {
	gorm.Model   `swaggerignore:"true"`
	Name         string      `json:"name" valid:"alphanum, maxstringlength(100)"`
	Publisher    string      `json:"publisher" valid:"alphanum, maxstringlength(100)"`
	PlayerNumber int         `json:"playerNumber" valid:"int, range(1|16)"`
	Tags         []Tag       `gorm:"many2many:boardgame_tags;" json:"tags,omitempty"`
	Categories   []Category  `gorm:"many2many:boardgame_categories;" json:"categories,omitempty"`
	Mechanisms   []Mechanism `gorm:"many2many:boardgame_mechanisms;" json:"mechanisms,omitempty"`
	Ratings      []Rating    `gorm:"many2many:boardgame_ratings;" json:"ratings,omitempty"`
	Expansions   []Boardgame `gorm:"foreignkey:BoardgameID" swaggerignore:"true" json:"expansions,omitempty"`
	BoardgameID  *uint       `swaggerignore:"true" json:"boardgame_id,omitempty"`
}

// Constructors
func NewBoardgame(name, publisher string, playerNumber int, tags []Tag, categories []Category, mechanisms []Mechanism) Boardgame {
	return Boardgame{
		Name:         name,
		Publisher:    publisher,
		PlayerNumber: playerNumber,
		Tags:         tags,
		Categories:   categories,
		Mechanisms:   mechanisms,
	}
}

func NewExpansion(name, publisher string, playerNumber int, tags []Tag, categories []Category, boardgameId *uint) Boardgame {
	return Boardgame{
		Name:         name,
		Publisher:    publisher,
		PlayerNumber: playerNumber,
		Tags:         tags,
		Categories:   categories,
		BoardgameID:  boardgameId,
	}
}

// Update
func (bg *Boardgame) UpdateBoardgame(boardgame Boardgame) {
	bg.Name = boardgame.GetName()

	bg.Publisher = boardgame.GetPublisher()

	bg.PlayerNumber = boardgame.GetPlayerNumber()

	bg.Tags = boardgame.GetTags()

	bg.Categories = boardgame.GetCategories()
}

// Existence functions
func (bg Boardgame) IsTags() bool {
	return len(bg.Tags) > 0
}

func (bg Boardgame) IsCategories() bool {
	return len(bg.Categories) > 0
}

func (bg Boardgame) IsMechanisms() bool {
	return len(bg.Mechanisms) > 0
}

func (bg Boardgame) IsExpansions() bool {
	return len(bg.Expansions) > 0
}

func (bg Boardgame) IsExpansion() bool {
	return bg.BoardgameID != nil
}

// Getters
func (bg Boardgame) GetId() *uint {
	return &bg.Model.ID
}

func (bg Boardgame) GetName() string {
	return bg.Name
}

func (bg Boardgame) GetPublisher() string {
	return bg.Publisher
}

func (bg Boardgame) GetPlayerNumber() int {
	return bg.PlayerNumber
}

func (bg Boardgame) GetTags() []Tag {
	return bg.Tags
}

func (bg Boardgame) GetCategories() []Category {
	return bg.Categories
}

func (bg Boardgame) GetMechanisms() []Mechanism {
	return bg.Mechanisms
}

func (bg Boardgame) GetExpansions() []Boardgame {
	return bg.Expansions
}

func (bg Boardgame) GetBoardgameID() *uint {
	return bg.BoardgameID
}

// Setters
func (bg *Boardgame) SetBoardgameID(id *uint) {
	bg.BoardgameID = id
}
