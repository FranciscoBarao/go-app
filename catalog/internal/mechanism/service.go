package mechanism

import (
	"context"
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package mechanism -destination service_mock.go . Database

// Database defines the persistence operations needed by the mechanism service.
type Database interface {
	CreateMechanism(ctx context.Context, m *Mechanism) error
	GetMechanismBySlug(ctx context.Context, slug string) (Mechanism, error)
	GetMechanismIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllMechanisms(ctx context.Context, query listopt.Query) ([]Mechanism, int, error)
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
	slugStr, err := svc.ensureUniqueSlug(ctx, req.Name)
	if err != nil {
		return Mechanism{}, err
	}

	m := Mechanism{
		Slug:  slugStr,
		Name:  req.Name,
		BggID: req.BggID,
	}
	if err := svc.db.CreateMechanism(ctx, &m); err != nil {
		return Mechanism{}, err
	}
	return m, nil
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active mechanism.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", middleware.NewError(http.StatusBadRequest, "name must contain slug-able characters")
	}

	_, err := svc.db.GetMechanismBySlug(ctx, slugStr)
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

// GetAll retrieves mechanisms with optional sort, filter, and pagination,
// returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, query listopt.Query) ([]Mechanism, int, error) {
	return svc.db.GetAllMechanisms(ctx, query)
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
