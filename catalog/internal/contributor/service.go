package contributor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/slug"
	"github.com/FranciscoBarao/catalog/internal/store"
)

//go:generate mockgen -package contributor -destination service_mock.go . Database

// Errors returned by the contributor service. Transport maps them to HTTP status codes.
var (
	ErrNotFound        = errors.New("record not found")
	ErrAlreadyExists   = errors.New("entry already exists")
	ErrValidation      = errors.New("validation failed")
	ErrStillReferenced = errors.New("still referenced")
)

// Database defines persistence operations for contributors.
type Database interface {
	CreateContributor(ctx context.Context, c *Contributor) error
	UpdateContributor(ctx context.Context, c *Contributor) error
	GetContributorBySlug(ctx context.Context, slug string) (Contributor, error)
	GetContributorIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllContributors(ctx context.Context, query listopt.Query) ([]Contributor, int, error)
	DeleteContributor(ctx context.Context, slug string, hard bool) error
}

// Service handles contributor business logic.
type Service struct {
	db Database
}

// NewService creates a new contributor Service.
func NewService(db Database) *Service {
	return &Service{db: db}
}

// Create builds a Contributor from the request, derives its slug, and persists it.
func (svc *Service) Create(ctx context.Context, req *CreateContributorRequest) (Contributor, error) {
	slugStr, err := svc.ensureUniqueSlug(ctx, req.Name)
	if err != nil {
		return Contributor{}, err
	}

	c := Contributor{
		Slug:  slugStr,
		Name:  req.Name,
		Bio:   req.Bio,
		BggID: req.BggID,
	}
	err = svc.db.CreateContributor(ctx, &c)
	switch {
	case err == nil:
		return c, nil
	case errors.Is(err, store.ErrAlreadyExists):
		return Contributor{}, fmt.Errorf("contributor %q: %w", slugStr, ErrAlreadyExists)
	default:
		return Contributor{}, fmt.Errorf("create contributor %q: internal error", slugStr)
	}
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active contributor.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", fmt.Errorf("name must contain slug-able characters: %w", ErrValidation)
	}

	_, err := svc.db.GetContributorBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return "", fmt.Errorf("slug %q already in use; disambiguate the name: %w", slugStr, ErrAlreadyExists)
	case errors.Is(err, store.ErrNotFound):
		return slugStr, nil
	default:
		return "", fmt.Errorf("get contributor %q: internal error", slugStr)
	}
}

// Update applies request fields onto the contributor identified by slug and persists it.
func (svc *Service) Update(ctx context.Context, req *UpdateContributorRequest, slugStr string) (Contributor, error) {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return Contributor{}, fmt.Errorf("name cannot be empty: %w", ErrValidation)
	}

	c, err := svc.Get(ctx, slugStr)
	if err != nil {
		return Contributor{}, err
	}
	req.ToContributor(&c)
	err = svc.db.UpdateContributor(ctx, &c)
	switch {
	case err == nil:
		return c, nil
	case errors.Is(err, store.ErrNotFound):
		return Contributor{}, fmt.Errorf("contributor %q: %w", slugStr, ErrNotFound)
	default:
		return Contributor{}, fmt.Errorf("update contributor %q: internal error", slugStr)
	}
}

// GetAll retrieves contributors with optional sort, filter, and pagination,
// returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, query listopt.Query) ([]Contributor, int, error) {
	contributors, total, err := svc.db.GetAllContributors(ctx, query)
	if err != nil {
		return nil, 0, errors.New("list contributors: internal error")
	}
	return contributors, total, nil
}

// Get retrieves a contributor by slug.
func (svc *Service) Get(ctx context.Context, slugStr string) (Contributor, error) {
	c, err := svc.db.GetContributorBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return c, nil
	case errors.Is(err, store.ErrNotFound):
		return Contributor{}, fmt.Errorf("contributor %q: %w", slugStr, ErrNotFound)
	default:
		return Contributor{}, fmt.Errorf("get contributor %q: internal error", slugStr)
	}
}

// GetIDBySlug validates a contributor exists and returns only its id.
func (svc *Service) GetIDBySlug(ctx context.Context, slugStr string) (uint, error) {
	id, err := svc.db.GetContributorIDBySlug(ctx, slugStr)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, store.ErrNotFound):
		return 0, fmt.Errorf("contributor %q: %w", slugStr, ErrNotFound)
	default:
		return 0, fmt.Errorf("get contributor %q: internal error", slugStr)
	}
}

// Delete removes a contributor by slug.
func (svc *Service) Delete(ctx context.Context, slugStr string, hard bool) error {
	err := svc.db.DeleteContributor(ctx, slugStr, hard)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, store.ErrNotFound):
		return fmt.Errorf("contributor %q: %w", slugStr, ErrNotFound)
	case errors.Is(err, store.ErrForeignKey):
		return fmt.Errorf("contributor %q: %w", slugStr, ErrStillReferenced)
	default:
		return fmt.Errorf("delete contributor %q: internal error", slugStr)
	}
}
