package repositories

import (
	"catalog/database"
	"catalog/middleware"
	"catalog/model"
	"errors"
)

type CategoryRepository struct {
	db *database.PostgresqlRepository
}

func NewCategoryRepository(instance *database.PostgresqlRepository) *CategoryRepository {
	return &CategoryRepository{
		db: instance,
	}
}

func (repo *CategoryRepository) Create(category *model.Category) error {

	return repo.db.Create(category)
}

func (repo *CategoryRepository) GetAll(sort string) ([]model.Category, error) {

	var categories []model.Category
	return categories, repo.db.Read(&categories, sort, "", "")
}

func (repo *CategoryRepository) Get(name string) (model.Category, error) {

	var category model.Category
	err := repo.db.Read(&category, "", "name = ?", name)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return category, middleware.NewError(mr.GetStatus(), "Category not found with name: "+name)
	}

	return category, err
}

func (repo *CategoryRepository) Delete(category *model.Category) error {

	return repo.db.Delete(category)
}
