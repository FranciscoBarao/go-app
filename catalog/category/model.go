package category

type Category struct {
	Name string `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
}

func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}

func (category *Category) UpdateCategory(name string) {
	if name != "" {
		category.Name = name
	}
}

// Getters
func (category Category) GetName() string {
	return category.Name
}
