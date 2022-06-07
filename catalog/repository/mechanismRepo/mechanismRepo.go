package mechanismRepo

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

func (repo *MechanismRepository) Create(mechanism model.Mechanism) error {

	err := repo.db.Create(&mechanism)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MechanismRepository) GetAll(sort string) ([]model.Mechanism, error) {

	var mechanisms []model.Mechanism
	err := repo.db.Read(&mechanisms, sort, "", "")
	if err != nil {
		return mechanisms, err
	}
	return mechanisms, nil
}

func (repo *MechanismRepository) Get(name string) (model.Mechanism, error) {

	var mechanism model.Mechanism
	err := repo.db.Read(&mechanism, "", "name = ?", name)
	if err != nil {
		return mechanism, err
	}
	return mechanism, nil
}

func (repo *MechanismRepository) Delete(mechanism model.Mechanism) error {

	err := repo.db.Delete(&mechanism)
	if err != nil {
		return err
	}
	return nil
}
