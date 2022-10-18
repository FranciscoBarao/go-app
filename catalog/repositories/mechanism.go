package repositories

import (
	"catalog/database"
	"catalog/model"
)

type MechanismRepository struct {
	db *database.PostgresqlRepository
}

func NewMechanismRepository(instance *database.PostgresqlRepository) *MechanismRepository {
	return &MechanismRepository{
		db: instance,
	}
}

func (repo *MechanismRepository) Create(mechanism *model.Mechanism) error {

	return repo.db.Create(&mechanism)
}

func (repo *MechanismRepository) GetAll(sort string) ([]model.Mechanism, error) {

	var mechanisms []model.Mechanism
	return mechanisms, repo.db.Read(&mechanisms, sort, "", "")
}

func (repo *MechanismRepository) Get(name string) (model.Mechanism, error) {

	var mechanism model.Mechanism
	return mechanism, repo.db.Read(&mechanism, "", "name = ?", name)
}

func (repo *MechanismRepository) Delete(mechanism model.Mechanism) error {

	return repo.db.Delete(&mechanism)
}
