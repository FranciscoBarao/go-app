package repositories

import (
	"catalog/database"
	"catalog/model"
)

type BoardgameRepository struct {
	db *database.PostgresqlRepository
}

func NewBoardgameRepository(instance *database.PostgresqlRepository) *BoardgameRepository {
	return &BoardgameRepository{
		db: instance,
	}
}

var omits = []string{"Tags.*", "Categories.*", "Mechanisms.*"}

func (repo *BoardgameRepository) Create(boardgame *model.Boardgame) error {

	return repo.db.Create(boardgame, omits...)
}

func (repo *BoardgameRepository) GetAll(sort, filterBody, filterValue string) ([]model.Boardgame, error) {

	var bg []model.Boardgame
	return bg, repo.db.Read(&bg, sort, filterBody, filterValue)
}

func (repo *BoardgameRepository) GetById(id string) (model.Boardgame, error) {

	var bg model.Boardgame
	return bg, repo.db.Read(&bg, "", "id = ?", id)
}

func (repo *BoardgameRepository) Update(boardgame model.Boardgame) error {

	if err := repo.db.Update(&boardgame, omits...); err != nil {
		return err
	}

	// Replace associations -> Easy fix? I dont like this approach -> Not modular
	tags := boardgame.GetTags()
	return repo.db.ReplaceAssociatons(&boardgame, "Tags", &tags)
}

func (repo *BoardgameRepository) DeleteById(boardgame model.Boardgame) error {

	return repo.db.Delete(&boardgame)
}