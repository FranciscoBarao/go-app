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

func (repo *BoardGameRepository) Create(boardgame model.BoardGame) error {

	err := repo.db.Create(&boardgame)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BoardGameRepository) GetAll() ([]model.BoardGame, error) {

	var bg []model.BoardGame
	err := repo.db.Read(&bg, "", "")
	if err != nil {
		return bg, err
	}
	return bg, nil
}

func (repo *BoardGameRepository) GetByName(name string) (model.BoardGame, error) {

	var bg model.BoardGame
	err := repo.db.Read(&bg, "name = ?", name)
	if err != nil {
		return bg, err
	}
	return bg, nil
}

func (repo *BoardGameRepository) GetById(id string) (model.BoardGame, error) {

	var bg model.BoardGame
	err := repo.db.Read(&bg, "id = ?", id)
	if err != nil {
		return bg, err
	}
	return bg, nil
}

func (repo *BoardGameRepository) Update(boardgame model.BoardGame) error {

	err := repo.db.Update(&boardgame)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoardGameRepository) DeleteById(boardgame model.BoardGame) error {

	err := repo.db.Delete(&boardgame)
	if err != nil {
		return err
	}
	return nil
}
