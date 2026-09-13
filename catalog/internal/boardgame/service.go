package boardgame

import (
	"context"
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package boardgame -destination service_mock.go . Database,CategoryService,MechanismService,ContributorService

// Database defines boardgame persistence operations.
type Database interface {
	CreateBoardgame(ctx context.Context, input CreateBoardgameDTO) (Boardgame, error)
	GetBoardgameByID(ctx context.Context, id uint) (Boardgame, error)
	GetBoardgameBySlug(ctx context.Context, slug string) (Boardgame, error)
	GetAllBoardgames(ctx context.Context, filter listopt.Params, includeDeleted bool) ([]Boardgame, int, error)
	UpdateBoardgame(ctx context.Context, bg *Boardgame) error
	UpdateBoardgameWithAssociations(ctx context.Context, bg *Boardgame, assoc UpdateAssociations) error
	DeleteBoardgame(ctx context.Context, id uint, hard bool) error
}

// CategoryService resolves a category id by slug.
type CategoryService interface {
	GetIDBySlug(ctx context.Context, slug string) (uint, error)
}

// MechanismService resolves a mechanism id by slug.
type MechanismService interface {
	GetIDBySlug(ctx context.Context, slug string) (uint, error)
}

// ContributorService resolves a contributor id by slug.
type ContributorService interface {
	GetIDBySlug(ctx context.Context, slug string) (uint, error)
}

// Service handles boardgame business logic.
type Service struct {
	db             Database
	categorySvc    CategoryService
	mechanismSvc   MechanismService
	contributorSvc ContributorService
}

// NewService creates a new boardgame Service.
func NewService(db Database, catSvc CategoryService, mechSvc MechanismService, contribSvc ContributorService) *Service {
	return &Service{
		db:             db,
		categorySvc:    catSvc,
		mechanismSvc:   mechSvc,
		contributorSvc: contribSvc,
	}
}

// Create validates the request, resolves associations, and persists a new boardgame.
func (svc *Service) Create(ctx context.Context, req *CreateBoardgameRequest, parentSlug string) (Boardgame, error) {
	logging.FromCtx(ctx).Debug().Str("parent_slug", parentSlug).Msg("creating boardgame")

	if err := validatePlayersOrder(req.MinPlayers, req.MaxPlayers); err != nil {
		return Boardgame{}, err
	}
	if err := validatePlayTimeOrder(req.MinPlayTime, req.MaxPlayTime); err != nil {
		return Boardgame{}, err
	}

	slugStr, err := svc.ensureUniqueSlug(ctx, req.Name)
	if err != nil {
		return Boardgame{}, err
	}

	parentID, err := svc.resolveParentID(ctx, parentSlug)
	if err != nil {
		return Boardgame{}, err
	}

	categoryIDs, err := svc.resolveCategoryRefs(ctx, req.Categories)
	if err != nil {
		return Boardgame{}, err
	}
	mechanismIDs, err := svc.resolveMechanismRefs(ctx, req.Mechanisms)
	if err != nil {
		return Boardgame{}, err
	}
	contributions, err := svc.resolveContributionInputs(ctx, req.Contributions)
	if err != nil {
		return Boardgame{}, err
	}

	dto := newCreateBoardgameDTO(req, slugStr, parentID, categoryIDs, mechanismIDs, contributions)
	return svc.db.CreateBoardgame(ctx, dto)
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active boardgame.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", middleware.NewError(http.StatusBadRequest, "name must contain slug-able characters")
	}

	_, err := svc.GetBySlug(ctx, slugStr)
	switch {
	case err == nil:
		logging.FromCtx(ctx).Error().Str("slug", slugStr).Msg("slug already in use")
		return "", middleware.NewError(http.StatusConflict, "slug '"+slugStr+"' already in use; disambiguate the name, e.g. add a year or edition")
	case middleware.IsNotFound(err):
		return slugStr, nil
	default:
		return "", err
	}
}

// GetAll retrieves boardgames with optional sort, filter, pagination, and deleted
// inclusion, returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, includeDeleted bool, params listopt.Params) ([]Boardgame, int, error) {
	return svc.db.GetAllBoardgames(ctx, params, includeDeleted)
}

// GetBySlug retrieves a boardgame by slug.
func (svc *Service) GetBySlug(ctx context.Context, slugStr string) (Boardgame, error) {
	return svc.db.GetBoardgameBySlug(ctx, slugStr)
}

// GetByID retrieves a boardgame by numeric id (for internal/service use).
func (svc *Service) GetByID(ctx context.Context, id uint) (Boardgame, error) {
	return svc.db.GetBoardgameByID(ctx, id)
}

// Update applies partial changes and persists.
func (svc *Service) Update(ctx context.Context, req *UpdateBoardgameRequest, slugStr string) error {
	logging.FromCtx(ctx).Debug().Str("slug", slugStr).Msg("updating boardgame")

	bg, err := svc.GetBySlug(ctx, slugStr)
	if err != nil {
		return err
	}

	req.ToBoardgame(&bg)

	if err := validatePlayersOrder(bg.MinPlayers, bg.MaxPlayers); err != nil {
		return err
	}
	if err := validatePlayTimeOrder(bg.MinPlayTime, bg.MaxPlayTime); err != nil {
		return err
	}

	assoc := UpdateAssociations{}
	if req.Categories != nil {
		categoryIDs, err := svc.resolveCategoryRefs(ctx, *req.Categories)
		if err != nil {
			return err
		}
		assoc.Categories = &categoryIDs
	}
	if req.Mechanisms != nil {
		mechanismIDs, err := svc.resolveMechanismRefs(ctx, *req.Mechanisms)
		if err != nil {
			return err
		}
		assoc.Mechanisms = &mechanismIDs
	}
	if req.Contributions != nil {
		contributions, err := svc.resolveContributionInputs(ctx, *req.Contributions)
		if err != nil {
			return err
		}
		assoc.Contributions = &contributions
	}

	return svc.db.UpdateBoardgameWithAssociations(ctx, &bg, assoc)
}

// DeleteBySlug soft-deletes or hard-deletes a boardgame.
func (svc *Service) DeleteBySlug(ctx context.Context, slugStr string, hard bool) error {
	bg, err := svc.GetBySlug(ctx, slugStr)
	if err != nil {
		return err
	}
	logging.FromCtx(ctx).Debug().Str("slug", slugStr).Bool("hard", hard).Msg("deleting boardgame")
	return svc.db.DeleteBoardgame(ctx, bg.ID, hard)
}

// Rate validates the boardgame exists and sets the username on the rating.
func (svc *Service) Rate(ctx context.Context, rating *Rating, slugStr string, username string) error {
	_, err := svc.GetBySlug(ctx, slugStr)
	if err != nil {
		return err
	}
	rating.Username = username
	// TODO - Redirect rating to another service
	return nil
}

func (svc *Service) resolveParentID(ctx context.Context, parentSlug string) (*uint, error) {
	if parentSlug == "" {
		return nil, nil
	}

	parent, err := svc.GetBySlug(ctx, parentSlug)
	if err != nil {
		return nil, err
	}

	if parent.IsExpansion() {
		logging.FromCtx(ctx).Error().Msg("an expansion cannot have other expansions")
		return nil, middleware.NewError(http.StatusConflict, "Expansion can't have expansions")
	}

	return &parent.ID, nil
}

func (svc *Service) resolveCategoryRefs(ctx context.Context, refs []CategoryRef) ([]uint, error) {
	if len(refs) == 0 {
		return nil, nil
	}
	ids := make([]uint, 0, len(refs))
	for _, ref := range refs {
		id, err := svc.categorySvc.GetIDBySlug(ctx, ref.Slug)
		if err != nil {
			logging.FromCtx(ctx).Error().Str("category", ref.Slug).Msg("category not found")
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (svc *Service) resolveMechanismRefs(ctx context.Context, refs []MechanismRef) ([]uint, error) {
	if len(refs) == 0 {
		return nil, nil
	}
	ids := make([]uint, 0, len(refs))
	for _, ref := range refs {
		id, err := svc.mechanismSvc.GetIDBySlug(ctx, ref.Slug)
		if err != nil {
			logging.FromCtx(ctx).Error().Str("mechanism", ref.Slug).Msg("mechanism not found")
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (svc *Service) resolveContributionInputs(ctx context.Context, inputs []contributor.ContributionInput) ([]ContributionDTO, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	contributions := make([]ContributionDTO, 0, len(inputs))
	for _, in := range inputs {
		if !in.Role.Valid() {
			logging.FromCtx(ctx).Error().Str("role", string(in.Role)).Msg("invalid contribution role")
			return nil, middleware.NewError(http.StatusBadRequest, "invalid contribution role")
		}
		id, err := svc.contributorSvc.GetIDBySlug(ctx, in.Slug)
		if err != nil {
			logging.FromCtx(ctx).Error().Str("contributor", in.Slug).Msg("contributor not found")
			return nil, err
		}
		contributions = append(contributions, ContributionDTO{
			ContributorID: id,
			Role:          in.Role,
			CreditOrder:   in.CreditOrder,
		})
	}
	return contributions, nil
}
