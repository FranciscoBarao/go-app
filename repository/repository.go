package repository

import (
	"go-app/database"
	"go-app/repository/boardgameRepo"
)

// Repositories contains all the repo structs
type Repositories struct {
	BoardGameRepository *boardgameRepo.BoardGameRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db *database.PostgresqlRepository) *Repositories {
	boardGameRepository := boardgameRepo.NewBoardGameRepository(db)

	return &Repositories{BoardGameRepository: boardGameRepository}
}
