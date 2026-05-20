package offer

import (
	"context"

	"github.com/FranciscoBarao/marketplace/internal/logging"
)

// Database defines the persistence operations needed by the offer service.
type Database interface {
	CreateOffer(ctx context.Context, offer *Offer) error
	GetOfferByID(ctx context.Context, uuid string) (Offer, error)
	GetOfferByIDAndUsername(ctx context.Context, uuid, username string) (Offer, error)
	GetAllOffers(ctx context.Context) ([]Offer, error)
	UpdateOffer(ctx context.Context, name string, price float64, uuid, username string) (Offer, error)
	DeleteOffer(ctx context.Context, uuid string) error
}

// Service holds the database dependency and implements offer business logic.
type Service struct {
	db Database
}

// NewService creates a new offer Service.
func NewService(db Database) *Service {
	return &Service{db: db}
}

// Create persists a new offer with the given username.
func (svc *Service) Create(ctx context.Context, offer *Offer, username string) error {
	logging.FromCtx(ctx).Debug().Str("username", username).Msg("creating offer")
	offer.Username = username
	return svc.db.CreateOffer(ctx, offer)
}

// GetAll retrieves all offers.
func (svc *Service) GetAll(ctx context.Context) ([]Offer, error) {
	return svc.db.GetAllOffers(ctx)
}

// Get retrieves a single offer by UUID.
func (svc *Service) Get(ctx context.Context, uuid string) (Offer, error) {
	return svc.db.GetOfferByID(ctx, uuid)
}

// Update applies changes and persists in a single query. Returns 404 if offer doesn't exist or isn't owned by username.
func (svc *Service) Update(ctx context.Context, req *UpdateOfferRequest, uuid, username string) (Offer, error) {
	logging.FromCtx(ctx).Debug().Str("uuid", uuid).Msg("updating offer")
	return svc.db.UpdateOffer(ctx, req.Name, req.Price, uuid, username)
}

// Delete removes an offer owned by the given username.
func (svc *Service) Delete(ctx context.Context, uuid, username string) error {
	logging.FromCtx(ctx).Debug().Str("uuid", uuid).Msg("deleting offer")

	offer, err := svc.db.GetOfferByIDAndUsername(ctx, uuid, username)
	if err != nil {
		return err
	}

	return svc.db.DeleteOffer(ctx, offer.UUID)
}
