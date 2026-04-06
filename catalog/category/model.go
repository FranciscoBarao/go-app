package category

// Category represents a board game category.
type Category struct {
	Name string `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
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
