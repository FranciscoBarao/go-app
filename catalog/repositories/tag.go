package repositories

import (
	"catalog/database"
	"catalog/model"
)

type TagRepository struct {
	db *database.PostgresqlRepository
}

func NewTagRepository(instance *database.PostgresqlRepository) *TagRepository {
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
	return tag, repo.db.Read(&tag, "", "name = ?", name)
}

func (repo *TagRepository) Delete(tag *model.Tag) error {

	return repo.db.Delete(tag)
}
