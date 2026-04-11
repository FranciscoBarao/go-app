package sql

// InsertCategory inserts a new category. Name is the primary key so duplicates will error.
const InsertCategory = `INSERT INTO categories (name) VALUES ($1)`

// SelectCategory retrieves a single category by name.
const SelectCategory = `SELECT name, created_at, updated_at
FROM categories
WHERE name = $1`

// SelectAllCategories retrieves all categories.
const SelectAllCategories = `SELECT name, created_at, updated_at
FROM categories`

// DeleteCategory removes a category by name.
const DeleteCategory = `DELETE FROM categories WHERE name = $1`
