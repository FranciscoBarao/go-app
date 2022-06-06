package repository

import (
	"go-app/database"
	"go-app/repository/boardgameRepo"
	"go-app/repository/categoryRepo"
	"go-app/repository/mechanismRepo"
	"go-app/repository/tagRepo"
)

// Repositories contains all the repo structs
type Repositories struct {
	BoardGameRepository *boardgameRepo.BoardGameRepository
	TagRepository       *tagRepo.TagRepository
	CategoryRepository  *categoryRepo.CategoryRepository
	MechanismRepository *mechanismRepo.MechanismRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	boardGameRepository := boardgameRepo.NewBoardGameRepository(db)
	tagRepository := tagRepo.NewTagRepository(db)
	categoryRepository := categoryRepo.NewCategoryRepository(db)
	mechanismRepository := mechanismRepo.NewMechanismRepository(db)

	return &Repositories{
		BoardGameRepository: boardGameRepository,
		TagRepository:       tagRepository,
		CategoryRepository:  categoryRepository,
		MechanismRepository: mechanismRepository,
	}
}
