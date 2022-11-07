package services

import "catalog/repositories"

// Repositories contains all the repo structs
type Services struct {
	BoardgameService *BoardgameService
	TagService       *TagService
	MechanismService *MechanismService
	CategoryService  *CategoryService
}

// InitRepositories should be called in main.go
func InitServices(repositories *repositories.Repositories) *Services {
	tagService := InitTagService(repositories.TagRepository)
	mechanismService := InitMechanismService(repositories.MechanismRepository)
	categoryService := InitCategoryService(repositories.CategoryRepository)
	boardgameService := InitBoardgameService(repositories.BoardgameRepository, tagService, categoryService, mechanismService)

	return &Services{
		BoardgameService: boardgameService,
		TagService:       tagService,
		MechanismService: mechanismService,
		CategoryService:  categoryService,
	}
}
