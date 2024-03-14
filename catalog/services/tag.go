package services

import (
	"catalog/model"
	"catalog/repositories"
)

type tagRepository interface {
	Create(tag *model.Tag) error
	GetAll(sort string) ([]model.Tag, error)
	Get(name string) (model.Tag, error)
	Delete(tag *model.Tag) error
}

type TagService struct {
	repo tagRepository
}

func InitTagService(tagRepo *repositories.TagRepository) *TagService {
	return &TagService{
		repo: tagRepo,
	}
}

func (svc *TagService) Create(tag *model.Tag) error {
	return svc.repo.Create(tag)
}

func (svc *TagService) GetAll(sort string) ([]model.Tag, error) {
	return svc.repo.GetAll(sort)
}

func (svc *TagService) Get(name string) (model.Tag, error) {
	return svc.repo.Get(name)
}

func (svc *TagService) Delete(name string) error {
	// Get category by name
	tag, err := svc.repo.Get(name)
	if err != nil {
		return err
	}

	// Delete by id
	return svc.repo.Delete(&tag)
}
