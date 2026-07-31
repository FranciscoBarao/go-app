package category

import (
	"context"
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package category -destination service_mock.go . Database

// Database defines the persistence operations needed by the category service.
type Database interface {
	CreateCategory(ctx context.Context, c *Category) error
	GetCategoryBySlug(ctx context.Context, slug string) (Category, error)
	GetCategoryIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllCategories(ctx context.Context, filter listopt.Params) ([]Category, int, error)
	DeleteCategory(ctx context.Context, slug string, hard bool) error
}

// Service handles category business logic.
type Service struct {
	db Database
}

// NewService creates a new category Service.
func NewService(db Database) *Service {
	return &Service{db: db}
}

// Create builds a Category from the request, derives its slug, and persists it.
func (svc *Service) Create(ctx context.Context, req *CreateCategoryRequest) (Category, error) {
	slugStr, err := svc.ensureUniqueSlug(ctx, req.Name)
	if err != nil {
		return Category{}, err
	}

	c := Category{
		Slug:  slugStr,
		Name:  req.Name,
		BggID: req.BggID,
	}
	if err := svc.db.CreateCategory(ctx, &c); err != nil {
		return Category{}, err
	}
	return c, nil
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active category.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", middleware.NewError(http.StatusBadRequest, "name must contain slug-able characters")
	}

	_, err := svc.db.GetCategoryBySlug(ctx, slugStr)
	switch {
	case err == nil:
		logging.FromCtx(ctx).Error().Str("slug", slugStr).Msg("slug already in use")
		return "", middleware.NewError(http.StatusConflict, "slug '"+slugStr+"' already in use; disambiguate the name")
	case middleware.IsNotFound(err):
		return slugStr, nil
	default:
		return "", err
	}
}

// GetAll retrieves categories with optional sort, filter, and pagination,
// returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, opts ...listopt.Option) ([]Category, int, error) {
	return svc.db.GetAllCategories(ctx, listopt.Apply(opts...))
}

// Get retrieves a Category by slug.
func (svc *Service) Get(ctx context.Context, slugStr string) (Category, error) {
	return svc.db.GetCategoryBySlug(ctx, slugStr)
}

// GetIDBySlug validates a Category exists and returns only its id.
func (svc *Service) GetIDBySlug(ctx context.Context, slugStr string) (uint, error) {
	return svc.db.GetCategoryIDBySlug(ctx, slugStr)
}

// Delete removes a Category by slug.
func (svc *Service) Delete(ctx context.Context, slugStr string, hard bool) error {
	return svc.db.DeleteCategory(ctx, slugStr, hard)
}
