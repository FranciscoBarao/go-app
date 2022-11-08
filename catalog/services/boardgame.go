package services

import (
	"catalog/middleware"
	"catalog/model"
	"catalog/repositories"
	"log"
	"net/http"
)

type boardgameRepository interface {
	Create(boardgame *model.Boardgame) error
	GetAll(sort, filterBody, filterValue string) ([]model.Boardgame, error)
	GetById(id string) (model.Boardgame, error)
	Update(boardgame model.Boardgame) error
	DeleteById(boardgame model.Boardgame) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type BoardgameService struct {
	repo         boardgameRepository
	tagSvc       *TagService
	categorySvc  *CategoryService
	mechanismSvc *MechanismService
}

// InitController initializes the boargame and the associations controller.
func InitBoardgameService(boardgameRepo *repositories.BoardgameRepository, tagService *TagService, categoryService *CategoryService, mechanismService *MechanismService) *BoardgameService {
	return &BoardgameService{
		repo:         boardgameRepo,
		tagSvc:       tagService,
		categorySvc:  categoryService,
		mechanismSvc: mechanismService,
	}
}

func (svc *BoardgameService) Create(boardgame *model.Boardgame, id string) error {

	// Check if Expansion -> Connect if needed
	if err := svc.connectBoardgameToExpansion(boardgame, id); err != nil {
		return err
	}

	// Check if Tags, Categories & Mechanisms exist
	if err := svc.validateAssociations(boardgame); err != nil {
		return err
	}

	return svc.repo.Create(boardgame)
}

func (svc *BoardgameService) GetAll(sort, filterBody, filterValue string) ([]model.Boardgame, error) {

	models, err := svc.repo.GetAll(sort, filterBody, filterValue)
	if err != nil {
		return models, err
	}

	return models, nil
}

func (svc *BoardgameService) GetById(id string) (model.Boardgame, error) {

	return svc.repo.GetById(id)
}

func (svc *BoardgameService) Update(input model.Boardgame, id string) error {

	// Check if Tags & Categories & Mechanisms exist
	if err := svc.validateAssociations(&input); err != nil {
		return err
	}

	// Get Boardgame by id
	boardgame, err := svc.repo.GetById(id)
	if err != nil {
		return err
	}

	// Updates Boardgame
	boardgame.UpdateBoardgame(input)

	return svc.repo.Update(boardgame)
}

func (svc *BoardgameService) DeleteById(id string) error {

	// Get Boardgame
	boardgame, err := svc.repo.GetById(id)
	if err != nil {
		return err
	}

	return svc.repo.DeleteById(boardgame)
}

// Function that checks if we are dealing with expansions and creates connection to boardgame parent
func (svc *BoardgameService) connectBoardgameToExpansion(boardgame *model.Boardgame, id string) error {
	if id != "" { // This is an expansion
		boardgameParent, err := svc.repo.GetById(id) // Get Parent BG
		if err != nil {
			return err
		}

		if boardgameParent.IsExpansion() {
			log.Println("Error -> An expansion cannot have other expansions")
			return middleware.NewError(http.StatusConflict, "Expansion can't have expansions")
		}

		boardgame.SetBoardgameID(boardgameParent.GetId()) // Set the Parents Id in the expansion
	}
	return nil
}

// Function that validates if tags, categories and mechanisms exist when boardgames are created
func (svc *BoardgameService) validateAssociations(boardgame *model.Boardgame) error {

	// Boardgame can contain Associations like Tags or Categories ->  We omit them which means that if they don't previously exist, the db returns an error -> Check if they exist before hand
	if boardgame.IsTags() {
		for _, tempTag := range boardgame.GetTags() {
			if _, err := svc.tagSvc.Get(tempTag.GetName()); err != nil { // Get tag by name
				return err // That tag does not exist -> Return Error
			}
		}
	}

	if boardgame.IsCategories() {
		for _, tempCategory := range boardgame.GetCategories() {
			if _, err := svc.categorySvc.Get(tempCategory.GetName()); err != nil { // Get category by name
				return err // That category does not exist -> Return Error
			}
		}
	}

	if boardgame.IsMechanisms() {
		for _, tempMechanism := range boardgame.GetMechanisms() {
			if _, err := svc.mechanismSvc.Get(tempMechanism.GetName()); err != nil { // Get mechanism by name
				return err // That mechanism does not exist -> Return Error
			}
		}
	}
	return nil
}
