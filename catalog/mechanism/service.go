package mechanism

import (
	"errors"

	"github.com/FranciscoBarao/catalog/database"
	"github.com/FranciscoBarao/catalog/middleware"
)

// Service merges the old MechanismRepository and Service into a single struct
// that holds a database.Database directly (no intermediate repository layer).
type Service struct {
	db database.Database
}

// NewService creates a new mechanism Service with the given database instance.
func NewService(db database.Database) *Service {
	return &Service{
		db: db,
	}
}

// Create persists a new Mechanism to the database.
func (svc *Service) Create(mechanism *Mechanism) error {
	return svc.db.Create(mechanism)
}

// GetAll retrieves all Mechanisms from the database, optionally sorted.
func (svc *Service) GetAll(sort string) ([]Mechanism, error) {
	var mechanisms []Mechanism
	return mechanisms, svc.db.Read(&mechanisms, sort, "", "")
}

// Get retrieves a single Mechanism by name. Returns a MalformedRequest error if not found.
func (svc *Service) Get(name string) (Mechanism, error) {
	var mechanism Mechanism
	err := svc.db.Read(&mechanism, "", "name = ?", name)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return mechanism, middleware.NewError(mr.GetStatus(), "Mechanism not found with name: "+name)
	}

	return mechanism, err
}

// Delete removes a Mechanism by name. It first retrieves the mechanism, then deletes it.
func (svc *Service) Delete(name string) error {
	mechanism, err := svc.Get(name)
	if err != nil {
		return err
	}

	return svc.db.Delete(&mechanism)
}
