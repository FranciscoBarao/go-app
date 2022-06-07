package categoryRepo

import (
	"catalog/database"
	"catalog/model"
)

type CategoryRepository struct {
	db *database.PostgresqlRepository
}

func NewCategoryRepository(instance *database.PostgresqlRepository) *CategoryRepository {
	return &CategoryRepository{
		db: instance,
	}
}

func (repo *CategoryRepository) Create(category model.Category) error {

	err := repo.db.Create(&category)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepository) GetAll(sort string) ([]model.Category, error) {

	var categories []model.Category
	err := repo.db.Read(&categories, sort, "", "")
	if err != nil {
		return categories, err
	}
	return categories, nil
}

func (repo *CategoryRepository) Get(name string) (model.Category, error) {

	var category model.Category
	err := repo.db.Read(&category, "", "name = ?", name)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (repo *CategoryRepository) Delete(category model.Category) error {

	err := repo.db.Delete(&category)
	if err != nil {
		return err
	}
	return nil
}
