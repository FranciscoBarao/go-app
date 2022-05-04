package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model `json:"-"`
	Name       string `gorm:"unique" json:"name"`
}

func NewTag(name string) Tag {
	return Tag{
		Name: name,
	}
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
