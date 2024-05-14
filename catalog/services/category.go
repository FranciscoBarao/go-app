package services

import (
	"github.com/FranciscoBarao/catalog/model"
	"github.com/FranciscoBarao/catalog/repositories"
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
	return svc.repo.GetAll(sort)
}

func (svc *CategoryService) Get(name string) (model.Category, error) {
	return svc.repo.Get(name)
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
