package model

import "gorm.io/gorm"

type Boardgame struct {
	gorm.Model   `swaggerignore:"true"`
	Name         string      `json:"name"`
	Publisher    string      `json:"publisher"`
	Price        float64     `json:"price"`
	PlayerNumber int         `json:"playerNumber"`
	Tags         []Tag       `gorm:"many2many:boardgame_tags;" json:"tags,omitempty"`
	Expansions   []Boardgame `gorm:"foreignkey:BoardgameID" swaggerignore:"true" json:"expansions,omitempty"`
	BoardgameID  *uint       `swaggerignore:"true" json:"boardgame_id,omitempty"`
}

func NewBoardgame(name, publisher string, price float64, playerNumber int, tags []Tag) Boardgame {
	return Boardgame{
		Name:         name,
		Publisher:    publisher,
		Price:        price,
		PlayerNumber: playerNumber,
		Tags:         tags,
	}
}

// Constructors
func (bg *Boardgame) UpdateBoardgame(name, publisher string, price float64, playerNumber int, tags []Tag) {
	bg.Name = name

	bg.Publisher = publisher

	bg.Price = price

	bg.PlayerNumber = playerNumber

	bg.Tags = tags
}

func NewExpansion(name, publisher string, price float64, playerNumber int, tags []Tag, boardgameId *uint) Boardgame {
	return Boardgame{
		Name:         name,
		Publisher:    publisher,
		Price:        price,
		PlayerNumber: playerNumber,
		Tags:         tags,
		BoardgameID:  boardgameId,
	}
}

// Existence functions
func (bg Boardgame) IsTags() bool {
	return len(bg.Tags) > 0
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

func (bg Boardgame) GetPrice() float64 {
	return bg.Price
}

func (bg Boardgame) GetPlayerNumber() int {
	return bg.PlayerNumber
}

func (bg Boardgame) GetTags() []Tag {
	return bg.Tags
}

func (bg Boardgame) GetExpansions() []Boardgame {
	return bg.Expansions
}

func (bg Boardgame) GetBoardgameID() *uint {
	return bg.BoardgameID
}

func (bg *Boardgame) SetBoardgameID(id *uint) {
	bg.BoardgameID = id
}
