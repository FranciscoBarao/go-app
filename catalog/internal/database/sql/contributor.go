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

// SoftDeleteContributor sets deleted_at on a contributor.
const SoftDeleteContributor = `UPDATE contributors SET deleted_at = NOW()
WHERE slug = $1 AND deleted_at IS NULL`

// HardDeleteContributor removes a contributor by slug.
const HardDeleteContributor = `DELETE FROM contributors WHERE slug = $1`

// SelectContributorBoardgames returns boardgames credited to a contributor, optionally filtered by role.
const SelectContributorBoardgames = `SELECT DISTINCT b.id, b.slug, b.created_at, b.updated_at, b.deleted_at, b.name, b.description,
    b.year_published, b.min_players, b.max_players, b.min_play_time, b.max_play_time, b.min_age, b.bgg_id, b.boardgame_id
FROM boardgames b
INNER JOIN boardgame_contributions bc ON bc.boardgame_id = b.id
INNER JOIN contributors c ON c.id = bc.contributor_id
WHERE c.slug = $1 AND b.deleted_at IS NULL AND c.deleted_at IS NULL
AND ($2::text IS NULL OR bc.role = $2)`
