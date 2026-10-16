package boardgame

import (
	"context"
	"errors"
	"fmt"

	"github.com/FranciscoBarao/catalog/internal/category"
	"github.com/FranciscoBarao/catalog/internal/contributor"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/mechanism"
	"github.com/FranciscoBarao/catalog/internal/slug"
	"github.com/FranciscoBarao/catalog/internal/store"
)

//go:generate mockgen -package boardgame -destination service_mock.go . Database,CategoryService,MechanismService,ContributorService

// Errors returned by the boardgame service. Transport maps them to HTTP status codes.
var (
	ErrNotFound          = errors.New("record not found")
	ErrAlreadyExists     = errors.New("entry already exists")
	ErrMissingReference  = errors.New("referenced record does not exist")
	ErrValidation        = errors.New("validation failed")
	ErrParentIsExpansion = errors.New("an expansion cannot be a parent")
)

// Database defines boardgame persistence operations.
type Database interface {
	CreateBoardgame(ctx context.Context, input CreateBoardgameDTO) (Boardgame, error)
	GetBoardgameByID(ctx context.Context, id uint) (Boardgame, error)
	GetBoardgameBySlug(ctx context.Context, slug string) (Boardgame, error)
	GetAllBoardgames(ctx context.Context, query listopt.Query, includeDeleted bool) ([]Boardgame, int, error)
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
	bg, err := svc.db.CreateBoardgame(ctx, dto)
	switch {
	case err == nil:
		return bg, nil
	case errors.Is(err, store.ErrAlreadyExists):
		return Boardgame{}, fmt.Errorf("boardgame %q: %w", slugStr, ErrAlreadyExists)
	case errors.Is(err, store.ErrForeignKey):
		return Boardgame{}, fmt.Errorf("boardgame %q: %w", slugStr, ErrMissingReference)
	default:
		return Boardgame{}, fmt.Errorf("create boardgame %q: internal error", slugStr)
	}
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active boardgame.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", fmt.Errorf("name must contain slug-able characters: %w", ErrValidation)
	}

	_, err := svc.GetBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return "", fmt.Errorf("slug %q already in use; disambiguate the name, e.g. add a year or edition: %w", slugStr, ErrAlreadyExists)
	case errors.Is(err, ErrNotFound):
		return slugStr, nil
	default:
		return "", err
	}
}

// GetAll retrieves boardgames with optional sort, filter, pagination, and deleted
// inclusion, returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, includeDeleted bool, query listopt.Query) ([]Boardgame, int, error) {
	boardgames, total, err := svc.db.GetAllBoardgames(ctx, query, includeDeleted)
	if err != nil {
		return nil, 0, errors.New("list boardgames: internal error")
	}
	return boardgames, total, nil
}

// GetBySlug retrieves a boardgame by slug.
func (svc *Service) GetBySlug(ctx context.Context, slugStr string) (Boardgame, error) {
	bg, err := svc.db.GetBoardgameBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return bg, nil
	case errors.Is(err, store.ErrNotFound):
		return Boardgame{}, fmt.Errorf("boardgame %q: %w", slugStr, ErrNotFound)
	default:
		return Boardgame{}, fmt.Errorf("get boardgame %q: internal error", slugStr)
	}
}

// GetByID retrieves a boardgame by numeric id (for internal/service use).
func (svc *Service) GetByID(ctx context.Context, id uint) (Boardgame, error) {
	bg, err := svc.db.GetBoardgameByID(ctx, id)
	switch {
	case err == nil:
		return bg, nil
	case errors.Is(err, store.ErrNotFound):
		return Boardgame{}, fmt.Errorf("boardgame %d: %w", id, ErrNotFound)
	default:
		return Boardgame{}, fmt.Errorf("get boardgame %d: internal error", id)
	}
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

	err = svc.db.UpdateBoardgameWithAssociations(ctx, &bg, assoc)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, store.ErrNotFound):
		return fmt.Errorf("boardgame %q: %w", slugStr, ErrNotFound)
	case errors.Is(err, store.ErrAlreadyExists):
		return fmt.Errorf("boardgame %q: %w", slugStr, ErrAlreadyExists)
	case errors.Is(err, store.ErrForeignKey):
		return fmt.Errorf("boardgame %q: %w", slugStr, ErrMissingReference)
	default:
		return fmt.Errorf("update boardgame %q: internal error", slugStr)
	}
}

// DeleteBySlug soft-deletes or hard-deletes a boardgame.
func (svc *Service) DeleteBySlug(ctx context.Context, slugStr string, hard bool) error {
	bg, err := svc.GetBySlug(ctx, slugStr)
	if err != nil {
		return err
	}
	logging.FromCtx(ctx).Debug().Str("slug", slugStr).Bool("hard", hard).Msg("deleting boardgame")
	err = svc.db.DeleteBoardgame(ctx, bg.ID, hard)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, store.ErrNotFound):
		return fmt.Errorf("boardgame %q: %w", slugStr, ErrNotFound)
	default:
		return fmt.Errorf("delete boardgame %q: internal error", slugStr)
	}
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
		return nil, fmt.Errorf("expansion %q: %w", parentSlug, ErrParentIsExpansion)
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
			if errors.Is(err, category.ErrNotFound) {
				return nil, fmt.Errorf("category %q: %w", ref.Slug, ErrNotFound)
			}
			return nil, fmt.Errorf("resolve category %q: internal error", ref.Slug)
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
			if errors.Is(err, mechanism.ErrNotFound) {
				return nil, fmt.Errorf("mechanism %q: %w", ref.Slug, ErrNotFound)
			}
			return nil, fmt.Errorf("resolve mechanism %q: internal error", ref.Slug)
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
			return nil, fmt.Errorf("contribution role %q: %w", in.Role, ErrValidation)
		}
		id, err := svc.contributorSvc.GetIDBySlug(ctx, in.Slug)
		if err != nil {
			if errors.Is(err, contributor.ErrNotFound) {
				return nil, fmt.Errorf("contributor %q: %w", in.Slug, ErrNotFound)
			}
			return nil, fmt.Errorf("resolve contributor %q: internal error", in.Slug)
		}
		contributions = append(contributions, ContributionDTO{
			ContributorID: id,
			Role:          in.Role,
			CreditOrder:   in.CreditOrder,
		})
	}
	return contributions, nil
}
