package sql

const contributorColumns = `id, slug, name, bio, bgg_id, created_at, updated_at, deleted_at`

// InsertContributor inserts a new contributor and returns generated fields.
const InsertContributor = `INSERT INTO contributors (slug, name, bio, bgg_id)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at`

// UpdateContributor updates mutable contributor fields and returns the refreshed updated_at.
const UpdateContributor = `UPDATE contributors
SET name = $1, bio = $2, bgg_id = $3
WHERE slug = $4 AND deleted_at IS NULL
RETURNING updated_at`

// SelectContributorBySlug retrieves a single active contributor by slug.
const SelectContributorBySlug = `SELECT ` + contributorColumns + `
FROM contributors
WHERE slug = $1 AND deleted_at IS NULL`

// SelectContributorIDBySlug retrieves only the id of an active contributor by slug.
const SelectContributorIDBySlug = `SELECT id FROM contributors
WHERE slug = $1 AND deleted_at IS NULL`

// SelectAllContributors retrieves all active contributors.
const SelectAllContributors = `SELECT ` + contributorColumns + `
FROM contributors
WHERE deleted_at IS NULL`

// CountContributors counts active contributors (before pagination).
const CountContributors = `SELECT COUNT(*) FROM contributors WHERE deleted_at IS NULL`

// SoftDeleteContributor sets deleted_at on a contributor.
const SoftDeleteContributor = `UPDATE contributors SET deleted_at = NOW()
WHERE slug = $1 AND deleted_at IS NULL`

// HardDeleteContributor removes a contributor by slug.
const HardDeleteContributor = `DELETE FROM contributors WHERE slug = $1`
