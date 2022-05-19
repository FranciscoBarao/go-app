package model

import "gorm.io/gorm"

type Boardgame struct {
	gorm.Model
	Name         string  `json:"name"`
	Dealer       string  `json:"dealer"`
	Price        float64 `json:"price"`
	PlayerNumber int     `json:"playerNumber"`
	Tags         []Tag   `gorm:"many2many:boardgame_tags;"`
}

func NewBoardgame(name, dealer string, price float64, playerNumber int, tags []Tag) Boardgame {
	return Boardgame{
		Name:         name,
		Dealer:       dealer,
		Price:        price,
		PlayerNumber: playerNumber,
		Tags:         tags,
	}
}

func (bg *Boardgame) UpdateBoardgame(name, dealer string, price float64, playerNumber int, tags []Tag) {
	bg.Name = name

	bg.Dealer = dealer

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

func (bg Boardgame) GetDealer() string {
	return bg.Dealer
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
