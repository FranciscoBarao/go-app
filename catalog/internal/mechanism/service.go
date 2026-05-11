package mechanism

import (
	"context"

	"github.com/FranciscoBarao/catalog/internal/query"
)

//go:generate mockgen -package mechanism -destination service_mock.go . Database

// Database defines the persistence operations needed by the mechanism service.
type Database interface {
	CreateMechanism(ctx context.Context, m *Mechanism) error
	GetMechanism(ctx context.Context, name string) (Mechanism, error)
	GetAllMechanisms(ctx context.Context, filter query.Filter) ([]Mechanism, error)
	DeleteMechanism(ctx context.Context, name string) error
}

// Service merges the old MechanismRepository and Service into a single struct
// that holds a Database directly (no intermediate repository layer).
type Service struct {
	db Database
}

// NewService creates a new mechanism Service with the given database instance.
func NewService(db Database) *Service {
	return &Service{
		db: db,
	}
}

// Create persists a new Mechanism to the database.
func (svc *Service) Create(ctx context.Context, mechanism *Mechanism) error {
	return svc.db.CreateMechanism(ctx, mechanism)
}

// GetAll retrieves all Mechanisms from the database, optionally sorted.
func (svc *Service) GetAll(ctx context.Context, opts ...query.Option) ([]Mechanism, error) {
	return svc.db.GetAllMechanisms(ctx, query.Apply(opts...))
}

// Get retrieves a single Mechanism by name.
func (svc *Service) Get(ctx context.Context, name string) (Mechanism, error) {
	return svc.db.GetMechanism(ctx, name)
}

// Delete removes a Mechanism by name.
func (svc *Service) Delete(ctx context.Context, name string) error {
	return svc.db.DeleteMechanism(ctx, name)
}
