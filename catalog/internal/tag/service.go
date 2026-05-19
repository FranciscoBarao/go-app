package tag

import (
	"context"

	"github.com/FranciscoBarao/catalog/internal/listopt"
)

//go:generate mockgen -package tag -destination service_mock.go . Database

// Database defines the persistence operations needed by the tag service.
type Database interface {
	CreateTag(ctx context.Context, t *Tag) error
	GetTag(ctx context.Context, name string) (Tag, error)
	GetAllTags(ctx context.Context, filter listopt.Params) ([]Tag, error)
	DeleteTag(ctx context.Context, name string) error
}

// Service merges the old TagRepository and Service into a single struct
// that holds a Database directly (no intermediate repository layer).
type Service struct {
	db Database
}

// NewService creates a new tag Service with the given database instance.
func NewService(db Database) *Service {
	return &Service{
		db: db,
	}
}

// Create persists a new Tag to the database.
func (svc *Service) Create(ctx context.Context, tag *Tag) error {
	return svc.db.CreateTag(ctx, tag)
}

// GetAll retrieves all Tags from the database, optionally sorted.
func (svc *Service) GetAll(ctx context.Context, opts ...listopt.Option) ([]Tag, error) {
	return svc.db.GetAllTags(ctx, listopt.Apply(opts...))
}

// Get retrieves a single Tag by name.
func (svc *Service) Get(ctx context.Context, name string) (Tag, error) {
	return svc.db.GetTag(ctx, name)
}

// Delete removes a Tag by name.
func (svc *Service) Delete(ctx context.Context, name string) error {
	return svc.db.DeleteTag(ctx, name)
}
