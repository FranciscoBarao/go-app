package mechanism

import (
	"context"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package mechanism -destination service_mock.go . Database

// Database defines the persistence operations needed by the mechanism service.
type Database interface {
	CreateMechanism(ctx context.Context, m *Mechanism) error
	GetMechanismBySlug(ctx context.Context, slug string) (Mechanism, error)
	GetMechanismIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllMechanisms(ctx context.Context, filter listopt.Params) ([]Mechanism, error)
	DeleteMechanism(ctx context.Context, slug string, hard bool) error
}

// Service handles mechanism business logic.
type Service struct {
	db Database
}

// NewService creates a new mechanism Service.
func NewService(db Database) *Service {
	return &Service{db: db}
}

// Create builds a Mechanism from the request, derives its slug, and persists it.
func (svc *Service) Create(ctx context.Context, req *CreateMechanismRequest) (Mechanism, error) {
	m := Mechanism{
		Slug:  slug.FromName(req.Name),
		Name:  req.Name,
		BggID: req.BggID,
	}
	if err := svc.db.CreateMechanism(ctx, &m); err != nil {
		return Mechanism{}, err
	}
	return m, nil
}

// GetAll retrieves all Mechanisms.
func (svc *Service) GetAll(ctx context.Context, opts ...listopt.Option) ([]Mechanism, error) {
	return svc.db.GetAllMechanisms(ctx, listopt.Apply(opts...))
}

// Get retrieves a Mechanism by slug.
func (svc *Service) Get(ctx context.Context, slugStr string) (Mechanism, error) {
	return svc.db.GetMechanismBySlug(ctx, slugStr)
}

// GetIDBySlug validates a Mechanism exists and returns only its id.
func (svc *Service) GetIDBySlug(ctx context.Context, slugStr string) (uint, error) {
	return svc.db.GetMechanismIDBySlug(ctx, slugStr)
}

// Delete removes a Mechanism by slug.
func (svc *Service) Delete(ctx context.Context, slugStr string, hard bool) error {
	return svc.db.DeleteMechanism(ctx, slugStr, hard)
}
