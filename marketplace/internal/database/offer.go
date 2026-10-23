package database

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"

	dbsql "github.com/FranciscoBarao/marketplace/internal/database/sql"
	"github.com/FranciscoBarao/marketplace/internal/middleware"
	"github.com/FranciscoBarao/marketplace/internal/offer"
)

// CreateOffer inserts a new offer and sets the generated UUID.
func (p *Postgres) CreateOffer(ctx context.Context, o *offer.Offer) error {
	err := p.pool.
		QueryRow(ctx, dbsql.InsertOffer, o.Username, o.Type, o.Name, o.Price).
		Scan(&o.UUID)
	return mapPgError(err)
}

// GetOfferByID retrieves a single offer by UUID.
func (p *Postgres) GetOfferByID(ctx context.Context, uuid string) (offer.Offer, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectOfferByID, uuid)
	if err != nil {
		return offer.Offer{}, mapPgError(err)
	}
	o, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[offer.Offer])
	return o, mapPgError(err)
}

// GetOfferByIDAndUsername retrieves a single offer by UUID and username.
func (p *Postgres) GetOfferByIDAndUsername(ctx context.Context, uuid, username string) (offer.Offer, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectOfferByIDAndUsername, uuid, username)
	if err != nil {
		return offer.Offer{}, mapPgError(err)
	}
	o, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[offer.Offer])
	return o, mapPgError(err)
}

// GetAllOffers retrieves all offers.
func (p *Postgres) GetAllOffers(ctx context.Context) ([]offer.Offer, error) {
	rows, err := p.pool.Query(ctx, dbsql.SelectAllOffers)
	if err != nil {
		return nil, mapPgError(err)
	}
	offers, err := pgx.CollectRows(rows, pgx.RowToStructByPos[offer.Offer])
	return offers, mapPgError(err)
}

// UpdateOffer updates the mutable fields of an offer owned by username and returns the updated row.
func (p *Postgres) UpdateOffer(ctx context.Context, name string, price float64, uuid, username string) (offer.Offer, error) {
	rows, err := p.pool.Query(ctx, dbsql.UpdateOffer, name, price, uuid, username)
	if err != nil {
		return offer.Offer{}, mapPgError(err)
	}
	o, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[offer.Offer])
	if err != nil {
		return offer.Offer{}, mapPgError(err)
	}
	return o, nil
}

// DeleteOffer removes an offer by UUID.
func (p *Postgres) DeleteOffer(ctx context.Context, uuid string) error {
	cmd, err := p.pool.Exec(ctx, dbsql.DeleteOffer, uuid)
	if err != nil {
		return mapPgError(err)
	}
	if cmd.RowsAffected() == 0 {
		return middleware.NewError(http.StatusNotFound, "record not found")
	}
	return nil
}
