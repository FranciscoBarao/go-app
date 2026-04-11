package sql

// InsertMechanism inserts a new mechanism. Name is the primary key so duplicates will error.
const InsertMechanism = `INSERT INTO mechanisms (name) VALUES ($1)`

// SelectMechanism retrieves a single mechanism by name.
const SelectMechanism = `SELECT name, created_at, updated_at
FROM mechanisms
WHERE name = $1`

// SelectAllMechanisms retrieves all mechanisms.
const SelectAllMechanisms = `SELECT name, created_at, updated_at
FROM mechanisms`

// DeleteMechanism removes a mechanism by name.
const DeleteMechanism = `DELETE FROM mechanisms WHERE name = $1`
