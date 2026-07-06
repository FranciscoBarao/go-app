package boardgame

import (
	"github.com/FranciscoBarao/catalog/internal/contributor"
)

// ContributionDTO carries the resolved contributor id and its credit metadata.
type ContributionDTO struct {
	ContributorID uint
	Role          contributor.Role
	CreditOrder   *int
}

// CreateBoardgameDTO is the write DTO passed from service to storage on create.
type CreateBoardgameDTO struct {
	Slug          string
	Name          string
	Description   *string
	YearPublished *int
	MinPlayers    int
	MaxPlayers    int
	MinPlayTime   *int
	MaxPlayTime   *int
	MinAge        *int
	BggID         *int
	ParentID      *uint
	Categories    []uint
	Mechanisms    []uint
	Contributions []ContributionDTO
}

func newCreateBoardgameDTO(
	req *CreateBoardgameRequest,
	slugStr string,
	parentID *uint,
	categoryIDs []uint,
	mechanismIDs []uint,
	contributions []ContributionDTO,
) CreateBoardgameDTO {
	return CreateBoardgameDTO{
		Slug:          slugStr,
		Name:          req.Name,
		Description:   req.Description,
		YearPublished: req.YearPublished,
		MinPlayers:    req.MinPlayers,
		MaxPlayers:    req.MaxPlayers,
		MinPlayTime:   req.MinPlayTime,
		MaxPlayTime:   req.MaxPlayTime,
		MinAge:        req.MinAge,
		BggID:         req.BggID,
		ParentID:      parentID,
		Categories:    categoryIDs,
		Mechanisms:    mechanismIDs,
		Contributions: contributions,
	}
}
