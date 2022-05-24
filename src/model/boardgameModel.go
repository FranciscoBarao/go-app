package model

import "gorm.io/gorm"

type Boardgame struct {
	gorm.Model   `swaggerignore:"true"`
	Name         string  `json:"name"`
	Publisher    string  `json:"publisher"`
	Price        float64 `json:"price"`
	PlayerNumber int     `json:"playerNumber"`
	Tags         []Tag   `gorm:"many2many:boardgame_tags;"`
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

func (bg *Boardgame) UpdateBoardgame(name, publisher string, price float64, playerNumber int, tags []Tag) {
	bg.Name = name

	bg.Publisher = publisher

	bg.Price = price

	bg.PlayerNumber = playerNumber

	bg.Tags = tags
}

func (bg Boardgame) IsTags() bool {
	return len(bg.Tags) > 0
}

// Getters
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
