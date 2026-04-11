package sql

// InsertBoardgame inserts a new boardgame and returns the generated id.
const InsertBoardgame = `INSERT INTO boardgames (name, publisher, player_number, boardgame_id)
VALUES ($1, $2, $3, $4)
RETURNING id`

// SelectBoardgameByID retrieves a single boardgame by id.
const SelectBoardgameByID = `SELECT id, created_at, updated_at, name, publisher, player_number, boardgame_id
FROM boardgames
WHERE id = $1`

// SelectAllBoardgames retrieves all boardgames.
const SelectAllBoardgames = `SELECT id, created_at, updated_at, name, publisher, player_number, boardgame_id
FROM boardgames`

// UpdateBoardgame updates a boardgame's mutable fields.
const UpdateBoardgame = `UPDATE boardgames
SET name = $1, publisher = $2, player_number = $3, boardgame_id = $4
WHERE id = $5`

// DeleteBoardgame removes a boardgame by id.
const DeleteBoardgame = `DELETE FROM boardgames WHERE id = $1`

// DeleteBoardgameTags removes all tag associations for a boardgame.
const DeleteBoardgameTags = `DELETE FROM boardgame_tags WHERE boardgame_id = $1`

// DeleteBoardgameCategories removes all category associations for a boardgame.
const DeleteBoardgameCategories = `DELETE FROM boardgame_categories WHERE boardgame_id = $1`

// DeleteBoardgameMechanisms removes all mechanism associations for a boardgame.
const DeleteBoardgameMechanisms = `DELETE FROM boardgame_mechanisms WHERE boardgame_id = $1`

// InsertBoardgameTag inserts a tag association for a boardgame.
const InsertBoardgameTag = `INSERT INTO boardgame_tags (boardgame_id, tag_name)
VALUES ($1, $2)`

// InsertBoardgameCategory inserts a category association for a boardgame.
const InsertBoardgameCategory = `INSERT INTO boardgame_categories (boardgame_id, category_name)
VALUES ($1, $2)`

// InsertBoardgameMechanism inserts a mechanism association for a boardgame.
const InsertBoardgameMechanism = `INSERT INTO boardgame_mechanisms (boardgame_id, mechanism_name)
VALUES ($1, $2)`

// SelectBoardgameTags retrieves all tags for a boardgame via join.
const SelectBoardgameTags = `SELECT t.name, t.created_at, t.updated_at
FROM tags t
INNER JOIN boardgame_tags bt ON bt.tag_name = t.name
WHERE bt.boardgame_id = $1`

// SelectBoardgameCategories retrieves all categories for a boardgame via join.
const SelectBoardgameCategories = `SELECT c.name, c.created_at, c.updated_at
FROM categories c
INNER JOIN boardgame_categories bc ON bc.category_name = c.name
WHERE bc.boardgame_id = $1`

// SelectBoardgameMechanisms retrieves all mechanisms for a boardgame via join.
const SelectBoardgameMechanisms = `SELECT m.name, m.created_at, m.updated_at
FROM mechanisms m
INNER JOIN boardgame_mechanisms bm ON bm.mechanism_name = m.name
WHERE bm.boardgame_id = $1`

// SelectBoardgameExpansions retrieves all expansion boardgames for a parent boardgame.
const SelectBoardgameExpansions = `SELECT id, created_at, updated_at, name, publisher, player_number, boardgame_id
FROM boardgames
WHERE boardgame_id = $1`
