package boardgame

import (
	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/tag"
)

// CreateBoardgameRequest is the request for POST /boardgame requests.
type CreateBoardgameRequest struct {
	Name         string                `json:"name" valid:"required,maxstringlength(100)"`
	Publisher    string                `json:"publisher" valid:"required,maxstringlength(100)"`
	PlayerNumber int                   `json:"playerNumber" valid:"required,range(1|16)"`
	Tags         []tag.Tag             `json:"tags,omitempty"`
	Categories   []category.Category   `json:"categories,omitempty"`
	Mechanisms   []mechanism.Mechanism `json:"mechanisms,omitempty"`
}

// UpdateBoardgameRequest is the request for PATCH /boardgame/{id} requests.
type UpdateBoardgameRequest struct {
	Name         *string                `json:"name,omitempty" valid:"optional,maxstringlength(100)"`
	Publisher    *string                `json:"publisher,omitempty" valid:"optional,maxstringlength(100)"`
	PlayerNumber *int                   `json:"playerNumber,omitempty" valid:"optional,range(1|16)"`
	Tags         *[]tag.Tag             `json:"tags,omitempty"`
	Categories   *[]category.Category   `json:"categories,omitempty"`
	Mechanisms   *[]mechanism.Mechanism `json:"mechanisms,omitempty"`
}

// NewBoardgame creates a Boardgame domain model from a CreateBoardgameRequest.
func NewBoardgame(req *CreateBoardgameRequest) *Boardgame {
	return &Boardgame{
		Name:         req.Name,
		Publisher:    req.Publisher,
		PlayerNumber: req.PlayerNumber,
		Tags:         req.Tags,
		Categories:   req.Categories,
		Mechanisms:   req.Mechanisms,
	}
}

// ToBoardgame applies non-nil fields from UpdateBoardgameRequest onto an existing Boardgame.
func (r *UpdateBoardgameRequest) ToBoardgame(bg *Boardgame) {
	if r.Name != nil {
		bg.Name = *r.Name
	}
	if r.Publisher != nil {
		bg.Publisher = *r.Publisher
	}
	if r.PlayerNumber != nil {
		bg.PlayerNumber = *r.PlayerNumber
	}
	if r.Tags != nil {
		bg.Tags = *r.Tags
	}
	if r.Categories != nil {
		bg.Categories = *r.Categories
	}
	if r.Mechanisms != nil {
		bg.Mechanisms = *r.Mechanisms
	}
}
