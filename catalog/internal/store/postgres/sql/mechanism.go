package sql

const mechanismColumns = `id, slug, name, bgg_id, created_at, updated_at, deleted_at`

// InsertMechanism inserts a new mechanism and returns generated fields.
const InsertMechanism = `INSERT INTO mechanisms (slug, name, bgg_id)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at`

// SelectMechanismBySlug retrieves a single active mechanism by slug.
const SelectMechanismBySlug = `SELECT ` + mechanismColumns + `
FROM mechanisms
WHERE slug = $1 AND deleted_at IS NULL`

// SelectMechanismIDBySlug retrieves only the id of an active mechanism by slug.
const SelectMechanismIDBySlug = `SELECT id FROM mechanisms
WHERE slug = $1 AND deleted_at IS NULL`

// SelectAllMechanisms retrieves all active mechanisms.
const SelectAllMechanisms = `SELECT ` + mechanismColumns + `
FROM mechanisms
WHERE deleted_at IS NULL`

// CountMechanisms counts active mechanisms (before pagination).
const CountMechanisms = `SELECT COUNT(*) FROM mechanisms WHERE deleted_at IS NULL`

// SoftDeleteMechanism sets deleted_at on a mechanism.
const SoftDeleteMechanism = `UPDATE mechanisms SET deleted_at = NOW()
WHERE slug = $1 AND deleted_at IS NULL`

// HardDeleteMechanism removes a mechanism by slug.
const HardDeleteMechanism = `DELETE FROM mechanisms WHERE slug = $1`
