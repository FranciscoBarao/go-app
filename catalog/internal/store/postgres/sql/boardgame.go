package sql

// Boardgame column list for SELECT queries (active rows use deleted_at IS NULL in WHERE).
const boardgameColumns = `id, slug, created_at, updated_at, deleted_at, name, description,
year_published, min_players, max_players, min_play_time, max_play_time, min_age, bgg_id, boardgame_id`

// InsertBoardgame inserts a new boardgame and returns the generated id.
const InsertBoardgame = `INSERT INTO boardgames (
    slug, name, description, year_published, min_players, max_players,
    min_play_time, max_play_time, min_age, bgg_id, boardgame_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id`

// SelectBoardgameByID retrieves a single active boardgame by id.
const SelectBoardgameByID = `SELECT ` + boardgameColumns + `
FROM boardgames
WHERE id = $1 AND deleted_at IS NULL`

// SelectBoardgameBySlug retrieves a single active boardgame by slug.
const SelectBoardgameBySlug = `SELECT ` + boardgameColumns + `
FROM boardgames
WHERE slug = $1 AND deleted_at IS NULL`

// SelectAllBoardgames retrieves active boardgames.
const SelectAllBoardgames = `SELECT ` + boardgameColumns + `
FROM boardgames
WHERE deleted_at IS NULL`

// SelectAllBoardgamesIncludingDeleted retrieves all boardgames, including soft-deleted ones.
const SelectAllBoardgamesIncludingDeleted = `SELECT ` + boardgameColumns + `
FROM boardgames
WHERE TRUE`

// CountBoardgames counts active boardgames (before pagination).
const CountBoardgames = `SELECT COUNT(*) FROM boardgames WHERE deleted_at IS NULL`

// CountBoardgamesIncludingDeleted counts all boardgames including soft-deleted (before pagination).
const CountBoardgamesIncludingDeleted = `SELECT COUNT(*) FROM boardgames WHERE TRUE`

// UpdateBoardgame updates a boardgame's mutable fields (slug is immutable).
const UpdateBoardgame = `UPDATE boardgames
SET name = $1, description = $2, year_published = $3,
    min_players = $4, max_players = $5, min_play_time = $6, max_play_time = $7,
    min_age = $8, bgg_id = $9, boardgame_id = $10
WHERE id = $11 AND deleted_at IS NULL`

// SoftDeleteBoardgame sets deleted_at on a boardgame.
const SoftDeleteBoardgame = `UPDATE boardgames SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL`

// SoftDeleteBoardgameExpansions soft-deletes expansions linked to a base game.
const SoftDeleteBoardgameExpansions = `UPDATE boardgames SET deleted_at = NOW()
WHERE boardgame_id = $1 AND deleted_at IS NULL`

// HardDeleteBoardgame removes a boardgame by id.
const HardDeleteBoardgame = `DELETE FROM boardgames WHERE id = $1`

// HardDeleteBoardgameExpansions removes expansions linked to a base game.
const HardDeleteBoardgameExpansions = `DELETE FROM boardgames WHERE boardgame_id = $1`

// DeleteBoardgameCategories removes all category associations for a boardgame.
const DeleteBoardgameCategories = `DELETE FROM boardgame_categories WHERE boardgame_id = $1`

// DeleteBoardgameMechanisms removes all mechanism associations for a boardgame.
const DeleteBoardgameMechanisms = `DELETE FROM boardgame_mechanisms WHERE boardgame_id = $1`

// DeleteBoardgameContributions removes all contributor associations for a boardgame.
const DeleteBoardgameContributions = `DELETE FROM boardgame_contributions WHERE boardgame_id = $1`

// InsertBoardgameCategory inserts a category association for a boardgame.
const InsertBoardgameCategory = `INSERT INTO boardgame_categories (boardgame_id, category_id)
VALUES ($1, $2)`

// InsertBoardgameMechanism inserts a mechanism association for a boardgame.
const InsertBoardgameMechanism = `INSERT INTO boardgame_mechanisms (boardgame_id, mechanism_id)
VALUES ($1, $2)`

// InsertBoardgameContribution inserts a contributor credit for a boardgame.
const InsertBoardgameContribution = `INSERT INTO boardgame_contributions (boardgame_id, contributor_id, role, credit_order)
VALUES ($1, $2, $3, $4)`

// SelectBoardgameCategories retrieves all categories for a boardgame via join.
const SelectBoardgameCategories = `SELECT c.id, c.slug, c.name, c.bgg_id, c.created_at, c.updated_at, c.deleted_at
FROM categories c
INNER JOIN boardgame_categories bc ON bc.category_id = c.id
WHERE bc.boardgame_id = $1 AND c.deleted_at IS NULL`

// SelectBoardgameMechanisms retrieves all mechanisms for a boardgame via join.
const SelectBoardgameMechanisms = `SELECT m.id, m.slug, m.name, m.bgg_id, m.created_at, m.updated_at, m.deleted_at
FROM mechanisms m
INNER JOIN boardgame_mechanisms bm ON bm.mechanism_id = m.id
WHERE bm.boardgame_id = $1 AND m.deleted_at IS NULL`

// SelectBoardgameContributions retrieves contributor credits for a boardgame.
const SelectBoardgameContributions = `SELECT c.id, c.slug, c.name, c.bio, c.bgg_id, c.created_at, c.updated_at, c.deleted_at,
       bc.role, bc.credit_order
FROM contributors c
INNER JOIN boardgame_contributions bc ON bc.contributor_id = c.id
WHERE bc.boardgame_id = $1 AND c.deleted_at IS NULL
ORDER BY bc.role, bc.credit_order NULLS LAST, c.name`

// SelectBoardgameExpansions retrieves expansion boardgames for a parent.
const SelectBoardgameExpansions = `SELECT ` + boardgameColumns + `
FROM boardgames
WHERE boardgame_id = $1 AND deleted_at IS NULL`
