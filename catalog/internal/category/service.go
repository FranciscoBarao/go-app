package category

import (
	"context"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package category -destination service_mock.go . Database

// Database defines the persistence operations needed by the category service.
type Database interface {
	CreateCategory(ctx context.Context, c *Category) error
	GetCategoryBySlug(ctx context.Context, slug string) (Category, error)
	GetCategoryIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllCategories(ctx context.Context, filter listopt.Params) ([]Category, error)
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
	c := Category{
		Slug:  slug.FromName(req.Name),
		Name:  req.Name,
		BggID: req.BggID,
	}
	if err := svc.db.CreateCategory(ctx, &c); err != nil {
		return Category{}, err
	}
	return c, nil
}

// GetAll retrieves all Categories.
func (svc *Service) GetAll(ctx context.Context, opts ...listopt.Option) ([]Category, error) {
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
