package tagRepo

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

func (repo *TagRepository) Create(tag model.Tag) error {

	err := repo.db.Create(&tag)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TagRepository) GetAll(sort string) ([]model.Tag, error) {

	var tags []model.Tag
	err := repo.db.Read(&tags, sort, "", "")
	if err != nil {
		return tags, err
	}
	return tags, nil
}

func (repo *TagRepository) Get(name string) (model.Tag, error) {

	var tag model.Tag
	err := repo.db.Read(&tag, "", "name = ?", name)
	if err != nil {
		return tag, err
	}
	return tag, nil
}

func (repo *TagRepository) Delete(tag model.Tag) error {

	err := repo.db.Delete(&tag)
	if err != nil {
		return err
	}
	return nil
}
