package tag

import (
	"errors"

	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/middleware"
)

// TagService merges the old TagRepository and TagService into a single struct
// that holds a database.Database directly (no intermediate repository layer).
type TagService struct {
	db database.Database
}

// NewTagService creates a new TagService with the given database instance.
func NewTagService(db database.Database) *TagService {
	return &TagService{
		db: db,
	}
}

// Create persists a new Tag to the database.
func (svc *TagService) Create(tag *Tag) error {
	return svc.db.Create(tag)
}

// GetAll retrieves all Tags from the database, optionally sorted.
func (svc *TagService) GetAll(sort string) ([]Tag, error) {
	var tags []Tag
	return tags, svc.db.Read(&tags, sort, "", "")
}

// Get retrieves a single Tag by name. Returns a MalformedRequest error if not found.
func (svc *TagService) Get(name string) (Tag, error) {
	var tag Tag
	err := svc.db.Read(&tag, "", "name = ?", name)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return tag, middleware.NewError(mr.GetStatus(), "Tag not found with name: "+name)
	}

	return tag, err
}

// Delete removes a Tag by name. It first retrieves the tag, then deletes it.
func (svc *TagService) Delete(name string) error {
	tag, err := svc.Get(name)
	if err != nil {
		return err
	}

	return svc.db.Delete(&tag)
}
