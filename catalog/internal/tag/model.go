package tag

import (
	"fmt"
	"time"
)

// Tag represents a board game tag.
type Tag struct {
	Name      string    `json:"name" valid:"alphanum, maxstringlength(30)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
