package category

import "time"

// Category represents a board game category.
type Category struct {
	Name      string    `json:"name" valid:"alphanum, maxstringlength(30)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewCategory creates a new Category with the given name.
func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}

// UpdateCategory updates the category name if non-empty.
func (category *Category) UpdateCategory(name string) {
	if name != "" {
		category.Name = name
	}
}
