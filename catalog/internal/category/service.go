package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/FranciscoBarao/catalog/internal/database/dberr"
	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package category -destination service_mock.go . Database

// Errors returned by the category service. Transport maps them to HTTP status codes.
var (
	ErrNotFound        = errors.New("record not found")
	ErrAlreadyExists   = errors.New("entry already exists")
	ErrValidation      = errors.New("validation failed")
	ErrStillReferenced = errors.New("still referenced")
)

// Database defines the persistence operations needed by the category service.
type Database interface {
	CreateCategory(ctx context.Context, c *Category) error
	GetCategoryBySlug(ctx context.Context, slug string) (Category, error)
	GetCategoryIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllCategories(ctx context.Context, query listopt.Query) ([]Category, int, error)
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
	err = svc.db.CreateCategory(ctx, &c)
	switch {
	case err == nil:
		return c, nil
	case errors.Is(err, dberr.ErrAlreadyExists):
		return Category{}, fmt.Errorf("category %q: %w", slugStr, ErrAlreadyExists)
	default:
		return Category{}, fmt.Errorf("create category %q: internal error", slugStr)
	}
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active category.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", fmt.Errorf("name must contain slug-able characters: %w", ErrValidation)
	}

	_, err := svc.db.GetCategoryBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return "", fmt.Errorf("slug %q already in use; disambiguate the name: %w", slugStr, ErrAlreadyExists)
	case errors.Is(err, dberr.ErrNotFound):
		return slugStr, nil
	default:
		return "", fmt.Errorf("get category %q: internal error", slugStr)
	}
}

// GetAll retrieves categories with optional sort, filter, and pagination,
// returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, query listopt.Query) ([]Category, int, error) {
	categories, total, err := svc.db.GetAllCategories(ctx, query)
	if err != nil {
		return nil, 0, errors.New("list categories: internal error")
	}
	return categories, total, nil
}

// Get retrieves a Category by slug.
func (svc *Service) Get(ctx context.Context, slugStr string) (Category, error) {
	c, err := svc.db.GetCategoryBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return c, nil
	case errors.Is(err, dberr.ErrNotFound):
		return Category{}, fmt.Errorf("category %q: %w", slugStr, ErrNotFound)
	default:
		return Category{}, fmt.Errorf("get category %q: internal error", slugStr)
	}
}

// GetIDBySlug validates a Category exists and returns only its id.
func (svc *Service) GetIDBySlug(ctx context.Context, slugStr string) (uint, error) {
	id, err := svc.db.GetCategoryIDBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, dberr.ErrNotFound):
		return 0, fmt.Errorf("category %q: %w", slugStr, ErrNotFound)
	default:
		return 0, fmt.Errorf("get category %q: internal error", slugStr)
	}
}

// Delete removes a Category by slug.
func (svc *Service) Delete(ctx context.Context, slugStr string, hard bool) error {
	err := svc.db.DeleteCategory(ctx, slugStr, hard)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, dberr.ErrNotFound):
		return fmt.Errorf("category %q: %w", slugStr, ErrNotFound)
	case errors.Is(err, dberr.ErrForeignKey):
		return fmt.Errorf("category %q: %w", slugStr, ErrStillReferenced)
	default:
		return fmt.Errorf("delete category %q: internal error", slugStr)
	}
}
