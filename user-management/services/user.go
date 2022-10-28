package services

import (
	"user-management/models"
	"user-management/repositories"
)

type userRepository interface {
	Register(user *models.User) error
	GetAll(sort string) ([]models.User, error)
	Get(username string) (models.User, error)
	Delete(user *models.User) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type UserService struct {
	repo userRepository
}

// InitController initializes the boargame and the associations controller.
func InitUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: userRepo,
	}
}

func (svc *UserService) Register(user *models.User) error {

	if err := user.HashPassword(user.GetPassword()); err != nil {
		return err
	}

	return svc.repo.Register(user)
}

func (svc *UserService) GetAll(sort string) ([]models.User, error) {

	users, err := svc.repo.GetAll(sort)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (svc *UserService) Login(username, password string) error {

	user, err := svc.repo.Get(username)
	if err != nil {
		return err
	}

	return user.CheckPassword(password)
}

func (svc *UserService) Delete(name string) error {

	user, err := svc.repo.Get(name)
	if err != nil {
		return err
	}

	// Delete by id
	return svc.repo.Delete(&user)
}
