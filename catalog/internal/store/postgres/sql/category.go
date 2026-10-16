package sql

const categoryColumns = `id, slug, name, bgg_id, created_at, updated_at, deleted_at`

// InsertCategory inserts a new category and returns generated fields.
const InsertCategory = `INSERT INTO categories (slug, name, bgg_id)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at`

// SelectCategoryBySlug retrieves a single active category by slug.
const SelectCategoryBySlug = `SELECT ` + categoryColumns + `
FROM categories
WHERE slug = $1 AND deleted_at IS NULL`

// SelectCategoryIDBySlug retrieves only the id of an active category by slug.
const SelectCategoryIDBySlug = `SELECT id FROM categories
WHERE slug = $1 AND deleted_at IS NULL`

// SelectAllCategories retrieves all active categories.
const SelectAllCategories = `SELECT ` + categoryColumns + `
FROM categories
WHERE deleted_at IS NULL`

// CountCategories counts active categories (before pagination).
const CountCategories = `SELECT COUNT(*) FROM categories WHERE deleted_at IS NULL`

// SoftDeleteCategory sets deleted_at on a category.
const SoftDeleteCategory = `UPDATE categories SET deleted_at = NOW()
WHERE slug = $1 AND deleted_at IS NULL`

// HardDeleteCategory removes a category by slug.
const HardDeleteCategory = `DELETE FROM categories WHERE slug = $1`
