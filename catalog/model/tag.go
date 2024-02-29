package model

import "fmt"

type Tag struct {
	Name       string      `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
	Boardgames []Boardgame `gorm:"many2many:boardgame_tags;" json:"-"`
}

func NewTag(name string) *Tag {
	return &Tag{
		Name: name,
	}
}

func (tag *Tag) String() string {
	return fmt.Sprintf("{ %s }", tag.Name)
}

func (tag *Tag) UpdateTag(name string) {
	if name != "" {
		tag.Name = name
	}
}

// Getters
func (tag Tag) GetName() string {
	return tag.Name
}
