package services

import (
	"catalog/model"
	"catalog/repositories"
)

type categoryRepository interface {
	Create(category *model.Category) error
	GetAll(sort string) ([]model.Category, error)
	Get(name string) (model.Category, error)
	Delete(category *model.Category) error
}

type CategoryService struct {
	repo categoryRepository
}

func InitCategoryService(tagRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: tagRepo,
	}
}

func (svc *CategoryService) Create(category *model.Category) error {

	return svc.repo.Create(category)
}

func (svc *CategoryService) GetAll(sort string) ([]model.Category, error) {

	categories, err := svc.repo.GetAll(sort)
	if err != nil {
		return categories, err
	}
	return categories, nil
}

func (svc *CategoryService) Get(name string) (model.Category, error) {

	category, err := svc.repo.Get(name)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (svc *CategoryService) Delete(name string) error {

	// Get category by name
	category, err := svc.repo.Get(name)
	if err != nil {
		return err
	}

	// Delete by id
	return svc.repo.Delete(&category)
}
