package contributor

import (
	"context"
	"net/http"
	"strings"

	"github.com/FranciscoBarao/catalog/internal/listopt"
	"github.com/FranciscoBarao/catalog/internal/logging"
	"github.com/FranciscoBarao/catalog/internal/middleware"
	"github.com/FranciscoBarao/catalog/internal/slug"
)

//go:generate mockgen -package contributor -destination service_mock.go . Database

// Database defines persistence operations for contributors.
type Database interface {
	CreateContributor(ctx context.Context, c *Contributor) error
	UpdateContributor(ctx context.Context, c *Contributor) error
	GetContributorBySlug(ctx context.Context, slug string) (Contributor, error)
	GetContributorIDBySlug(ctx context.Context, slug string) (uint, error)
	GetAllContributors(ctx context.Context, filter listopt.Params) ([]Contributor, int, error)
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
	if err := svc.db.CreateContributor(ctx, &c); err != nil {
		return Contributor{}, err
	}
	return c, nil
}

// ensureUniqueSlug derives a slug from name, rejects empty slugs, and fails on collision with an active contributor.
func (svc *Service) ensureUniqueSlug(ctx context.Context, name string) (string, error) {
	slugStr := slug.FromName(name)
	if slugStr == "" {
		return "", middleware.NewError(http.StatusBadRequest, "name must contain slug-able characters")
	}

	_, err := svc.db.GetContributorBySlug(ctx, slugStr)
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

// Update applies request fields onto the contributor identified by slug and persists it.
func (svc *Service) Update(ctx context.Context, req *UpdateContributorRequest, slugStr string) (Contributor, error) {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return Contributor{}, middleware.NewError(http.StatusBadRequest, "name cannot be empty")
	}

	c, err := svc.db.GetContributorBySlug(ctx, slugStr)
	if err != nil {
		return Contributor{}, err
	}
	req.ToContributor(&c)
	if err := svc.db.UpdateContributor(ctx, &c); err != nil {
		return Contributor{}, err
	}
	return c, nil
}

// GetAll retrieves contributors with optional sort, filter, and pagination,
// returning the page of results and the total count of matching rows.
func (svc *Service) GetAll(ctx context.Context, opts ...listopt.Option) ([]Contributor, int, error) {
	return svc.db.GetAllContributors(ctx, listopt.Apply(opts...))
}

// Get retrieves a contributor by slug.
func (svc *Service) Get(ctx context.Context, slugStr string) (Contributor, error) {
	return svc.db.GetContributorBySlug(ctx, slugStr)
}

// GetIDBySlug validates a contributor exists and returns only its id.
func (svc *Service) GetIDBySlug(ctx context.Context, slugStr string) (uint, error) {
	return svc.db.GetContributorIDBySlug(ctx, slugStr)
}

// Delete removes a contributor by slug.
func (svc *Service) Delete(ctx context.Context, slugStr string, hard bool) error {
	return svc.db.DeleteContributor(ctx, slugStr, hard)
}
