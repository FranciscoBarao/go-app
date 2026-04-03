package category

import (
	"errors"

	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/middleware"
)

// Service merges the old CategoryRepository and Service into a single struct
// that holds a database.Database directly (no intermediate repository layer).
type Service struct {
	db database.Database
}

// NewService creates a new category Service with the given database instance.
func NewService(db database.Database) *Service {
	return &Service{
		db: db,
	}
}

// Create persists a new Category to the database.
func (svc *Service) Create(category *Category) error {
	return svc.db.Create(category)
}

// GetAll retrieves all Categories from the database, optionally sorted.
func (svc *Service) GetAll(sort string) ([]Category, error) {
	var categories []Category
	return categories, svc.db.Read(&categories, sort, "", "")
}

// Get retrieves a single Category by name. Returns a MalformedRequest error if not found.
func (svc *Service) Get(name string) (Category, error) {
	var category Category
	err := svc.db.Read(&category, "", "name = ?", name)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return category, middleware.NewError(mr.GetStatus(), "Category not found with name: "+name)
	}

	return category, err
}

// Delete removes a Category by name. It first retrieves the category, then deletes it.
func (svc *Service) Delete(name string) error {
	category, err := svc.Get(name)
	if err != nil {
		return err
	}

	return svc.db.Delete(&category)
}
