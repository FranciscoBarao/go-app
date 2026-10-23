package sql

// InsertOffer inserts a new offer and returns the generated UUID.
const InsertOffer = `INSERT INTO offers (username, type, name, price) VALUES ($1, $2, $3, $4) RETURNING uuid`

// SelectOfferByID retrieves a single offer by UUID.
const SelectOfferByID = `SELECT uuid, created_at, updated_at, username, type, name, price FROM offers WHERE uuid = $1`

// SelectOfferByIDAndUsername retrieves a single offer by UUID and username.
const SelectOfferByIDAndUsername = `SELECT uuid, created_at, updated_at, username, type, name, price FROM offers WHERE uuid = $1 AND username = $2`

// SelectAllOffers retrieves all offers.
const SelectAllOffers = `SELECT uuid, created_at, updated_at, username, type, name, price FROM offers`

// UpdateOffer updates the mutable fields of an offer and returns the full row.
const UpdateOffer = `UPDATE offers SET name = $1, price = $2 WHERE uuid = $3 AND username = $4 RETURNING uuid, created_at, updated_at, username, type, name, price`

// DeleteOffer removes an offer by UUID.
const DeleteOffer = `DELETE FROM offers WHERE uuid = $1`
