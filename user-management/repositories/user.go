package repositories

import (
	"user-management/database"
	"user-management/models"
)

type UserRepository struct {
	db *database.PostgresqlRepository
}

func NewUserRepository(instance *database.PostgresqlRepository) *UserRepository {
	return &UserRepository{
		db: instance,
	}
}

func (repo *UserRepository) Register(user *models.User) error {

	return repo.db.Create(user)
}

func (repo *UserRepository) GetAll(sort string) ([]models.User, error) {

	var users []models.User
	return users, repo.db.Read(&users, sort, "", "")
}

func (repo *UserRepository) Get(name string) (models.User, error) {

	var user models.User
	return user, repo.db.Read(&user, "", "name = ?", name)
}

func (repo *UserRepository) Delete(user *models.User) error {

	return repo.db.Delete(user)
}
