package tag

import "fmt"

// Tag represents a board game tag.
type Tag struct {
	Name string `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
}

// NewTag creates a new Tag with the given name.
func NewTag(name string) *Tag {
	return &Tag{
		Name: name,
	}
}

func (tag *Tag) String() string {
	return fmt.Sprintf("{ %s }", tag.Name)
}

// UpdateTag updates the tag name if non-empty.
func (tag *Tag) UpdateTag(name string) {
	if name != "" {
		tag.Name = name
	}
}
