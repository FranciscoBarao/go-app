package services

import (
	"catalog/model"
	"catalog/repositories"
)

type mechanismRepository interface {
	Create(mechanism *model.Mechanism) error
	GetAll(sort string) ([]model.Mechanism, error)
	Get(name string) (model.Mechanism, error)
	Delete(mechanism model.Mechanism) error
}

type MechanismService struct {
	repo mechanismRepository
}

func InitMechanismService(mechanismRepo *repositories.MechanismRepository) *MechanismService {
	return &MechanismService{
		repo: mechanismRepo,
	}
}

func (svc *MechanismService) Create(mechanism *model.Mechanism) error {

	return svc.repo.Create(mechanism)
}

func (svc *MechanismService) GetAll(sort string) ([]model.Mechanism, error) {

	mechanisms, err := svc.repo.GetAll(sort)
	if err != nil {
		return mechanisms, err
	}
	return mechanisms, nil
}

func (svc *MechanismService) Get(name string) (model.Mechanism, error) {

	mechanism, err := svc.repo.Get(name)
	if err != nil {
		return mechanism, err
	}
	return mechanism, nil
}

func (svc *MechanismService) Delete(name string) error {

	// Get Mechanism by name
	mechanism, err := svc.repo.Get(name)
	if err != nil {
		return err
	}

	return svc.repo.Delete(mechanism)
}
