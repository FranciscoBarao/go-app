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

	tags, err := svc.repo.GetAll(sort)
	if err != nil {
		return tags, err
	}
	return tags, nil
}

func (svc *TagService) Get(name string) (model.Tag, error) {

	tag, err := svc.repo.Get(name)
	if err != nil {
		return tag, err
	}
	return tag, nil
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
