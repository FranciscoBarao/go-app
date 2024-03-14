package model

type Category struct {
	Name       string      `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
	Boardgames []Boardgame `gorm:"many2many:boardgame_categories;" json:"-"`
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
