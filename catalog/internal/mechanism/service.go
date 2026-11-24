package mechanism

import (
	"context"
	"errors"
	"fmt"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/slug"
	"github.com/FranciscoBarao/catalog/internal/store"
)

//go:generate mockgen -package mechanism -destination service_mock.go . Database

// Errors returned by the mechanism service. Transport maps them to HTTP status codes.
var (
	ErrNotFound        = errors.New("record not found")
	ErrAlreadyExists   = errors.New("entry already exists")
	ErrValidation      = errors.New("validation failed")
	ErrStillReferenced = errors.New("still referenced")
)

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
	err = svc.db.CreateMechanism(ctx, &m)
	switch {
	case err == nil:
		logging.FromCtx(ctx).Info().Str("slug", m.Slug).Msg("mechanism created")
		return m, nil
	case errors.Is(err, store.ErrAlreadyExists):
		return Mechanism{}, fmt.Errorf("mechanism %q: %w", slugStr, ErrAlreadyExists)
	default:
		return Mechanism{}, fmt.Errorf("create mechanism %q: internal error", slugStr)
	}
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active mechanism.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", fmt.Errorf("name must contain slug-able characters: %w", ErrValidation)
	}

	_, err := svc.db.GetMechanismBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return "", fmt.Errorf("slug %q already in use; disambiguate the name: %w", slugStr, ErrAlreadyExists)
	case errors.Is(err, store.ErrNotFound):
		return slugStr, nil
	default:
		return "", fmt.Errorf("get mechanism %q: internal error", slugStr)
	}
}

// GetAll retrieves mechanisms with optional sort, filter, and pagination,
// returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, query listopt.Query) ([]Mechanism, int, error) {
	mechanisms, total, err := svc.db.GetAllMechanisms(ctx, query)
	if err != nil {
		return nil, 0, errors.New("list mechanisms: internal error")
	}
	return mechanisms, total, nil
}

// Get retrieves a Mechanism by slug.
func (svc *Service) Get(ctx context.Context, slugStr string) (Mechanism, error) {
	m, err := svc.db.GetMechanismBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return m, nil
	case errors.Is(err, store.ErrNotFound):
		return Mechanism{}, fmt.Errorf("mechanism %q: %w", slugStr, ErrNotFound)
	default:
		return Mechanism{}, fmt.Errorf("get mechanism %q: internal error", slugStr)
	}
}

// GetIDBySlug validates a Mechanism exists and returns only its id.
func (svc *Service) GetIDBySlug(ctx context.Context, slugStr string) (uint, error) {
	id, err := svc.db.GetMechanismIDBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, store.ErrNotFound):
		return 0, fmt.Errorf("mechanism %q: %w", slugStr, ErrNotFound)
	default:
		return 0, fmt.Errorf("get mechanism %q: internal error", slugStr)
	}
}

// Delete removes a Mechanism by slug.
func (svc *Service) Delete(ctx context.Context, slugStr string, hard bool) error {
	err := svc.db.DeleteMechanism(ctx, slugStr, hard)
	switch {
	case err == nil:
		logging.FromCtx(ctx).Info().Str("slug", slugStr).Bool("hard", hard).Msg("mechanism deleted")
		return nil
	case errors.Is(err, store.ErrNotFound):
		return fmt.Errorf("mechanism %q: %w", slugStr, ErrNotFound)
	case errors.Is(err, store.ErrForeignKey):
		return fmt.Errorf("mechanism %q: %w", slugStr, ErrStillReferenced)
	default:
		return fmt.Errorf("delete mechanism %q: internal error", slugStr)
	}
}
