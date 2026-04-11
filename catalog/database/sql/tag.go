package sql

// InsertTag inserts a new tag. Name is the primary key so duplicates will error.
const InsertTag = `INSERT INTO tags (name) VALUES ($1)`

// SelectTag retrieves a single tag by name.
const SelectTag = `SELECT name, created_at, updated_at
FROM tags
WHERE name = $1`

// SelectAllTags retrieves all tags.
const SelectAllTags = `SELECT name, created_at, updated_at
FROM tags`

// DeleteTag removes a tag by name.
const DeleteTag = `DELETE FROM tags WHERE name = $1`
