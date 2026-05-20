package boardgame

import (
	"context"
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/tag"
)

//go:generate mockgen -package boardgame -destination service_mock.go . Database,TagService,CategoryService,MechanismService

// Database defines the interface for database operations needed by the boardgame service.
type Database interface {
	CreateBoardgame(ctx context.Context, bg *Boardgame) error
	GetBoardgameByID(ctx context.Context, id uint) (Boardgame, error)
	GetAllBoardgames(ctx context.Context, filter listopt.Params) ([]Boardgame, error)
	UpdateBoardgame(ctx context.Context, bg *Boardgame) error
	UpdateBoardgameWithAssociations(ctx context.Context, bg *Boardgame, assoc UpdateAssociations) error
	DeleteBoardgame(ctx context.Context, id uint) error
}

// TagService is a local interface for retrieving tags by name.
type TagService interface {
	Get(ctx context.Context, name string) (tag.Tag, error)
}

// CategoryService is a local interface for retrieving categories by name.
type CategoryService interface {
	Get(ctx context.Context, name string) (category.Category, error)
}

// MechanismService is a local interface for retrieving mechanisms by name.
type MechanismService interface {
	Get(ctx context.Context, name string) (mechanism.Mechanism, error)
}

// Service merges the old BoardgameRepository and Service into a single struct
// that holds a Database directly and uses local getter interfaces for association validation.
type Service struct {
	db           Database
	tagSvc       TagService
	categorySvc  CategoryService
	mechanismSvc MechanismService
}

// NewService creates a new boardgame Service with the given dependencies.
func NewService(db Database, tagSvc TagService, catSvc CategoryService, mechSvc MechanismService) *Service {
	return &Service{
		db:           db,
		tagSvc:       tagSvc,
		categorySvc:  catSvc,
		mechanismSvc: mechSvc,
	}
}

// Create validates associations, connects expansions if needed, and persists a new Boardgame.
func (svc *Service) Create(ctx context.Context, boardgame *Boardgame, id uint) error {
	logging.FromCtx(ctx).Debug().Uint("parent_id", id).Msg("creating boardgame")

	// Check if Expansion -> Connect if needed
	if err := svc.connectBoardgameToExpansion(ctx, boardgame, id); err != nil {
		return err
	}

	// Check if Tags, Categories & Mechanisms exist
	if err := svc.validateAssociations(ctx, boardgame); err != nil {
		return err
	}

	return svc.db.CreateBoardgame(ctx, boardgame)
}

// GetAll retrieves all Boardgames from the database with optional sort and filter.
func (svc *Service) GetAll(ctx context.Context, opts ...listopt.Option) ([]Boardgame, error) {
	return svc.db.GetAllBoardgames(ctx, listopt.Apply(opts...))
}

// GetByID retrieves a single Boardgame by its ID.
func (svc *Service) GetByID(ctx context.Context, id uint) (Boardgame, error) {
	return svc.db.GetBoardgameByID(ctx, id)
}

// Update fetches the existing boardgame, applies partial changes from the request, validates associations, and persists.
func (svc *Service) Update(ctx context.Context, req *UpdateBoardgameRequest, id uint) error {
	logging.FromCtx(ctx).Debug().Uint("id", id).Msg("updating boardgame")
	// Get Boardgame by id
	boardgame, err := svc.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply partial update
	req.ToBoardgame(&boardgame)

	// Check if Tags & Categories & Mechanisms exist
	if err := svc.validateAssociations(ctx, &boardgame); err != nil {
		return err
	}

	return svc.db.UpdateBoardgameWithAssociations(ctx, &boardgame, UpdateAssociations{
		Tags:       req.Tags,
		Categories: req.Categories,
		Mechanisms: req.Mechanisms,
	})
}

// DeleteByID deletes a boardgame by its ID.
func (svc *Service) DeleteByID(ctx context.Context, id uint) error {
	logging.FromCtx(ctx).Debug().Uint("id", id).Msg("deleting boardgame")
	return svc.db.DeleteBoardgame(ctx, id)
}

// Rate validates the boardgame exists and sets the username on the rating.
func (svc *Service) Rate(ctx context.Context, rating *Rating, id uint, username string) error {
	// Check if boardgame exists
	_, err := svc.GetByID(ctx, id)
	if err != nil {
		return err
	}

	rating.Username = username

	// TODO - Redirect rating to another service

	return nil
}

// connectBoardgameToExpansion checks if we are dealing with expansions and creates connection to boardgame parent.
func (svc *Service) connectBoardgameToExpansion(ctx context.Context, boardgame *Boardgame, id uint) error {
	if id == 0 { // This is an expansion
		return nil
	}

	boardgameParent, err := svc.GetByID(ctx, id) // Get Parent BG
	if err != nil {
		return err
	}

	if boardgameParent.IsExpansion() {
		logging.FromCtx(ctx).Error().Msg("an expansion cannot have other expansions")
		return middleware.NewError(http.StatusConflict, "Expansion can't have expansions")
	}

	boardgame.BoardgameID = &boardgameParent.ID // Set the Parents Id in the expansion
	return nil
}

// validateAssociations validates if tags, categories and mechanisms exist when boardgames are created.
func (svc *Service) validateAssociations(ctx context.Context, boardgame *Boardgame) error {
	if boardgame.HasTags() {
		for _, tempTag := range boardgame.Tags {
			if _, err := svc.tagSvc.Get(ctx, tempTag.Name); err != nil {
				logging.FromCtx(ctx).Error().Str("tag", tempTag.Name).Msg("tag not found")
				return err
			}
		}
	}

	if boardgame.HasCategories() {
		for _, tempCategory := range boardgame.Categories {
			if _, err := svc.categorySvc.Get(ctx, tempCategory.Name); err != nil {
				logging.FromCtx(ctx).Error().Str("category", tempCategory.Name).Msg("category not found")
				return err
			}
		}
	}

	if boardgame.HasMechanisms() {
		for _, tempMechanism := range boardgame.Mechanisms {
			if _, err := svc.mechanismSvc.Get(ctx, tempMechanism.Name); err != nil {
				logging.FromCtx(ctx).Error().Str("mechanism", tempMechanism.Name).Msg("mechanism not found")
				return err
			}
		}
	}
	return nil
}
