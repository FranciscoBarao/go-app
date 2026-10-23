package boardgame

import (
	"github.com/FranciscoBarao/catalog/internal/contributor"
)

// CategoryRef references a category by slug.
type CategoryRef struct {
	Slug string `json:"slug" valid:"required,maxstringlength(100)"`
}

// MechanismRef references a mechanism by slug.
type MechanismRef struct {
	Slug string `json:"slug" valid:"required,maxstringlength(100)"`
}

// CreateBoardgameRequest is the request for POST /boardgame requests.
type CreateBoardgameRequest struct {
	Name          string                          `json:"name" valid:"required,maxstringlength(120)"`
	Description   string                          `json:"description" valid:"optional"`
	YearPublished int                             `json:"year_published" valid:"optional,range(1900|2100)"`
	MinPlayers    int                             `json:"min_players" valid:"required,range(1|16)"`
	MaxPlayers    int                             `json:"max_players" valid:"required,range(1|16)"`
	MinPlayTime   int                             `json:"min_play_time" valid:"optional,range(1|9999)"`
	MaxPlayTime   int                             `json:"max_play_time" valid:"optional,range(1|9999)"`
	MinAge        *int                            `json:"min_age,omitempty" valid:"optional,range(0|99)"`
	BggID         *int                            `json:"bgg_id,omitempty"`
	Categories    []CategoryRef                   `json:"categories,omitempty"`
	Mechanisms    []MechanismRef                  `json:"mechanisms,omitempty"`
	Contributions []contributor.ContributionInput `json:"contributions,omitempty"`
}

// UpdateBoardgameRequest is the request for PATCH /boardgame/{slug} requests.
type UpdateBoardgameRequest struct {
	Name          *string                          `json:"name,omitempty" valid:"optional,maxstringlength(120)"`
	Description   *string                          `json:"description,omitempty"`
	YearPublished *int                             `json:"year_published,omitempty" valid:"optional,range(1900|2100)"`
	MinPlayers    *int                             `json:"min_players,omitempty" valid:"optional,range(1|16)"`
	MaxPlayers    *int                             `json:"max_players,omitempty" valid:"optional,range(1|16)"`
	MinPlayTime   *int                             `json:"min_play_time,omitempty" valid:"optional,range(1|9999)"`
	MaxPlayTime   *int                             `json:"max_play_time,omitempty" valid:"optional,range(1|9999)"`
	MinAge        *int                             `json:"min_age,omitempty" valid:"optional,range(0|99)"`
	BggID         *int                             `json:"bgg_id,omitempty"`
	Categories    *[]CategoryRef                   `json:"categories,omitempty"`
	Mechanisms    *[]MechanismRef                  `json:"mechanisms,omitempty"`
	Contributions *[]contributor.ContributionInput `json:"contributions,omitempty"`
}

// ToBoardgame applies non-nil scalar fields from UpdateBoardgameRequest onto an existing Boardgame.
func (r *UpdateBoardgameRequest) ToBoardgame(bg *Boardgame) {
	if r.Name != nil {
		bg.Name = *r.Name
	}
	if r.Description != nil {
		bg.Description = *r.Description
	}
	if r.YearPublished != nil {
		bg.YearPublished = *r.YearPublished
	}
	if r.MinPlayers != nil {
		bg.MinPlayers = *r.MinPlayers
	}
	if r.MaxPlayers != nil {
		bg.MaxPlayers = *r.MaxPlayers
	}
	if r.MinPlayTime != nil {
		bg.MinPlayTime = *r.MinPlayTime
	}
	if r.MaxPlayTime != nil {
		bg.MaxPlayTime = *r.MaxPlayTime
	}
	if r.MinAge != nil {
		bg.MinAge = r.MinAge
	}
	if r.BggID != nil {
		bg.BggID = r.BggID
	}
}

// UpdateAssociations signals which association tables should be replaced during an update.
type UpdateAssociations struct {
	Categories    *[]uint
	Mechanisms    *[]uint
	Contributions *[]ContributionDTO
}
