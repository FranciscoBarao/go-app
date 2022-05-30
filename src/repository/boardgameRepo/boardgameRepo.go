package boardgameRepo

import (
	"go-app/database"
	"go-app/model"
)

type BoardGameRepository struct {
	db *database.PostgresqlRepository
}

func NewBoardGameRepository(instance *database.PostgresqlRepository) *BoardGameRepository {
	return &BoardGameRepository{
		db: instance,
	}
}

var omits = []string{"Tags.*"}

func (repo *BoardGameRepository) Create(boardgame model.Boardgame) error {

	err := repo.db.Create(&boardgame, omits...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BoardGameRepository) GetAll(filterBody, filterValue string) ([]model.Boardgame, error) {

	var bg []model.Boardgame
	err := repo.db.Read(&bg, filterBody, filterValue)
	if err != nil {
		return bg, err
	}
	return bg, nil
}

func (repo *BoardGameRepository) GetByName(name string) (model.Boardgame, error) {

	var bg model.Boardgame
	err := repo.db.Read(&bg, "name = ?", name)
	if err != nil {
		return bg, err
	}
	return bg, nil
}

func (repo *BoardGameRepository) GetById(id string) (model.Boardgame, error) {

	var bg model.Boardgame
	err := repo.db.Read(&bg, "id = ?", id)
	if err != nil {
		return bg, err
	}
	return bg, nil
}

func (repo *BoardGameRepository) Update(boardgame model.Boardgame) error {

	err := repo.db.Update(&boardgame, omits...)
	if err != nil {
		return err
	}

	// Replace associations -> Easy fix? I dont like this approach -> Not modular
	tags := boardgame.GetTags()
	err = repo.db.ReplaceAssociatons(&boardgame, "Tags", &tags)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoardGameRepository) DeleteById(boardgame model.Boardgame) error {

	err := repo.db.Delete(&boardgame)
	if err != nil {
		return err
	}
	return nil
}
