package repositories

import (
	"catalog/database"
)

// Repositories contains all the repo structs
type Repositories struct {
	BoardgameRepository *BoardgameRepository
	TagRepository       *TagRepository
	CategoryRepository  *CategoryRepository
	MechanismRepository *MechanismRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	boardgameRepository := NewBoardgameRepository(db)
	tagRepository := NewTagRepository(db)
	categoryRepository := NewCategoryRepository(db)
	mechanismRepository := NewMechanismRepository(db)

	return &Repositories{
		BoardgameRepository: boardgameRepository,
		TagRepository:       tagRepository,
		CategoryRepository:  categoryRepository,
		MechanismRepository: mechanismRepository,
	}
}
