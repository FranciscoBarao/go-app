package controllers

import "catalog/repositories"

// Controllers contains all the controllers
type Controllers struct {
	BoardgameController *BoardgameController
	TagController       *TagController
	CategoryController  *CategoryController
	MechanismController *MechanismController
}

// InitControllers returns a new Controllers
func InitControllers(repositories *repositories.Repositories) *Controllers {
	return &Controllers{
		BoardgameController: InitBoardgameController(repositories.BoardgameRepository, repositories.TagRepository, repositories.CategoryRepository, repositories.MechanismRepository),
		TagController:       InitTagController(repositories.TagRepository),
		CategoryController:  InitCategoryController(repositories.CategoryRepository),
		MechanismController: InitMechanismController(repositories.MechanismRepository),
	}
}
