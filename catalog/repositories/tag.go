package repositories

import (
	"catalog/middleware"
	"catalog/model"
	"errors"
)

type TagRepository struct {
	db Database
}

func NewTagRepository(instance Database) *TagRepository {
	return &TagRepository{
		db: instance,
	}
}

func (repo *TagRepository) Create(tag *model.Tag) error {
	return repo.db.Create(tag)
}

func (repo *TagRepository) GetAll(sort string) ([]model.Tag, error) {
	var tags []model.Tag
	return tags, repo.db.Read(&tags, sort, "", "")
}

func (repo *TagRepository) Get(name string) (model.Tag, error) {
	var tag model.Tag
	err := repo.db.Read(&tag, "", "name = ?", name)

	var mr *middleware.MalformedRequest
	if err != nil && errors.As(err, &mr) {
		return tag, middleware.NewError(mr.GetStatus(), "Tag not found with name: "+name)
	}

	return tag, err
}

func (repo *TagRepository) Delete(tag *model.Tag) error {
	return repo.db.Delete(tag)
}
