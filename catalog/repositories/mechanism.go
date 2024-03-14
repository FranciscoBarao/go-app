package repositories

import (
	"catalog/middleware"
	"catalog/model"
	"errors"
)

type MechanismRepository struct {
	db Database
}

func NewMechanismRepository(instance Database) *MechanismRepository {
	return &MechanismRepository{
		db: instance,
	}
}

func (repo *MechanismRepository) Create(mechanism *model.Mechanism) error {

	return repo.db.Create(mechanism)
}

func (repo *MechanismRepository) GetAll(sort string) ([]model.Mechanism, error) {

	var mechanisms []model.Mechanism
	return mechanisms, repo.db.Read(&mechanisms, sort, "", "")
}

func (repo *MechanismRepository) Get(name string) (model.Mechanism, error) {

	var mechanism model.Mechanism
	err := repo.db.Read(&mechanism, "", "name = ?", name)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return mechanism, middleware.NewError(mr.GetStatus(), "Mechanism not found with name: "+name)
	}

	return mechanism, err
}

func (repo *MechanismRepository) Delete(mechanism model.Mechanism) error {

	return repo.db.Delete(&mechanism)
}
