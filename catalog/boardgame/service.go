package boardgame

import (
	"context"
	"errors"
	"net/http"

	"github.com/FranciscoBarao/catalog/category"
	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/mechanism"
	"github.com/FranciscoBarao/catalog/middleware"
	"github.com/FranciscoBarao/catalog/tag"
)

// TagGetter is a local interface for retrieving tags by name.
type TagGetter interface {
	Get(name string) (tag.Tag, error)
}

// CategoryGetter is a local interface for retrieving categories by name.
type CategoryGetter interface {
	Get(name string) (category.Category, error)
}

// MechanismGetter is a local interface for retrieving mechanisms by name.
type MechanismGetter interface {
	Get(name string) (mechanism.Mechanism, error)
}

// BoardgameService merges the old BoardgameRepository and BoardgameService into a single struct
// that holds a database.Database directly and uses local getter interfaces for association validation.
type BoardgameService struct {
	db           database.Database
	tagSvc       TagGetter
	categorySvc  CategoryGetter
	mechanismSvc MechanismGetter
}

// NewBoardgameService creates a new BoardgameService with the given dependencies.
func NewBoardgameService(db database.Database, tagSvc TagGetter, catSvc CategoryGetter, mechSvc MechanismGetter) *BoardgameService {
	return &BoardgameService{
		db:           db,
		tagSvc:       tagSvc,
		categorySvc:  catSvc,
		mechanismSvc: mechSvc,
	}
}

// Create validates associations, connects expansions if needed, and persists a new Boardgame.
func (svc *BoardgameService) Create(boardgame *Boardgame, id string) error {
	// Check if Expansion -> Connect if needed
	if err := svc.connectBoardgameToExpansion(boardgame, id); err != nil {
		return err
	}

	// Check if Tags, Categories & Mechanisms exist
	if err := svc.validateAssociations(boardgame); err != nil {
		return err
	}

	return svc.db.Create(boardgame)
}

// GetAll retrieves all Boardgames from the database with optional sort and filter.
func (svc *BoardgameService) GetAll(sort, filterBody, filterValue string) ([]Boardgame, error) {
	var bg []Boardgame
	return bg, svc.db.Read(&bg, sort, filterBody, filterValue)
}

// GetById retrieves a single Boardgame by its ID.
func (svc *BoardgameService) GetById(id string) (Boardgame, error) {
	var bg Boardgame
	err := svc.db.Read(&bg, "", "id = ?", id)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return bg, middleware.NewError(mr.GetStatus(), "Boardgame not found with id: "+id)
	}

	return bg, err
}

// Update validates associations, fetches the existing boardgame, applies changes, and persists.
func (svc *BoardgameService) Update(input *Boardgame, id string) error {
	// Check if Tags & Categories & Mechanisms exist
	if err := svc.validateAssociations(input); err != nil {
		return err
	}

	// Get Boardgame by id
	boardgame, err := svc.GetById(id)
	if err != nil {
		return err
	}

	// Updates Boardgame
	boardgame.UpdateBoardgame(input)

	if err := svc.db.Update(&boardgame); err != nil {
		return err
	}

	// Replace associations -> Easy fix? I dont like this approach -> Not modular
	tags := boardgame.GetTags()
	return svc.db.ReplaceAssociatons(&boardgame, "Tags", &tags)
}

// DeleteById fetches a boardgame by ID and deletes it.
func (svc *BoardgameService) DeleteById(id string) error {
	// Get Boardgame
	boardgame, err := svc.GetById(id)
	if err != nil {
		return err
	}

	return svc.db.Delete(&boardgame)
}

// Rate validates the boardgame exists and sets the username on the rating.
func (svc *BoardgameService) Rate(rating *Rating, id, username string) error {
	// Check if boardgame exists
	_, err := svc.GetById(id)
	if err != nil {
		return err
	}

	rating.SetUsername(username)

	// TODO - Redirect rating to another service

	return nil
}

// connectBoardgameToExpansion checks if we are dealing with expansions and creates connection to boardgame parent.
func (svc *BoardgameService) connectBoardgameToExpansion(boardgame *Boardgame, id string) error {
	if id == "" { // This is an expansion
		return nil
	}

	boardgameParent, err := svc.GetById(id) // Get Parent BG
	if err != nil {
		return err
	}

	if boardgameParent.IsExpansion() {
		middleware.FromCtx(context.Background()).Error().Msg("an expansion cannot have other expansions")
		return middleware.NewError(http.StatusConflict, "Expansion can't have expansions")
	}

	boardgame.SetBoardgameID(boardgameParent.GetId()) // Set the Parents Id in the expansion
	return nil
}

// validateAssociations validates if tags, categories and mechanisms exist when boardgames are created.
func (svc *BoardgameService) validateAssociations(boardgame *Boardgame) error {
	// Boardgame can contain Associations like Tags or Categories ->  We omit them which means that if they don't previously exist, the db returns an error -> Check if they exist before hand
	if boardgame.HasTags() {
		for _, tempTag := range boardgame.GetTags() {
			if _, err := svc.tagSvc.Get(tempTag.GetName()); err != nil { // Get tag by name
				return err // That tag does not exist -> Return Error
			}
		}
	}

	if boardgame.HasCategories() {
		for _, tempCategory := range boardgame.GetCategories() {
			if _, err := svc.categorySvc.Get(tempCategory.GetName()); err != nil { // Get category by name
				return err // That category does not exist -> Return Error
			}
		}
	}

	if boardgame.HasMechanisms() {
		for _, tempMechanism := range boardgame.GetMechanisms() {
			if _, err := svc.mechanismSvc.Get(tempMechanism.GetName()); err != nil { // Get mechanism by name
				return err // That mechanism does not exist -> Return Error
			}
		}
	}
	return nil
}
