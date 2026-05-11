package category

import (
	"context"

	"github.com/FranciscoBarao/catalog/internal/query"
)

//go:generate mockgen -package category -destination service_mock.go . Database

// Database defines the persistence operations needed by the category service.
type Database interface {
	CreateCategory(ctx context.Context, c *Category) error
	GetCategory(ctx context.Context, name string) (Category, error)
	GetAllCategories(ctx context.Context, filter query.Filter) ([]Category, error)
	DeleteCategory(ctx context.Context, name string) error
}

// Service merges the old CategoryRepository and Service into a single struct
// that holds a Database directly (no intermediate repository layer).
type Service struct {
	db Database
}

// NewService creates a new category Service with the given database instance.
func NewService(db Database) *Service {
	return &Service{
		db: db,
	}
}

// Create persists a new Category to the database.
func (svc *Service) Create(ctx context.Context, category *Category) error {
	return svc.db.CreateCategory(ctx, category)
}

// GetAll retrieves all Categories from the database, optionally sorted.
func (svc *Service) GetAll(ctx context.Context, opts ...query.Option) ([]Category, error) {
	return svc.db.GetAllCategories(ctx, query.Apply(opts...))
}

// Get retrieves a single Category by name.
func (svc *Service) Get(ctx context.Context, name string) (Category, error) {
	return svc.db.GetCategory(ctx, name)
}

// Delete removes a Category by name.
func (svc *Service) Delete(ctx context.Context, name string) error {
	return svc.db.DeleteCategory(ctx, name)
}
