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

type UserService struct {
	repo userRepository
}

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

	return svc.repo.GetAll(sort)
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
