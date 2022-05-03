package model

import "gorm.io/gorm"

type BoardGame struct {
	gorm.Model
	Name         string  `json:"name"`
	Dealer       string  `json:"dealer"`
	Price        float64 `json:"price"`
	PlayerNumber int     `json:"playerNumber"`
}

func NewBoardGame(name, dealer string, price float64, playerNumber int) BoardGame {
	return BoardGame{
		Name:         name,
		Dealer:       dealer,
		Price:        price,
		PlayerNumber: playerNumber,
	}
}

func (bg *BoardGame) UpdateBoardGame(name, dealer string, price float64, playerNumber int) {
	if name != "" {
		bg.Name = name
	}
	if dealer != "" {
		bg.Dealer = dealer
	}
	if price > 0 {
		bg.Price = price
	}

	if playerNumber > 0 {
		bg.PlayerNumber = playerNumber
	}
}

// Getters
func (bg BoardGame) GetName() string {
	return bg.Name
}

func (bg BoardGame) GetDealer() string {
	return bg.Dealer
}

func (bg BoardGame) GetPrice() float64 {
	return bg.Price
}

func (bg BoardGame) GetPlayerNumber() int {
	return bg.PlayerNumber
}
