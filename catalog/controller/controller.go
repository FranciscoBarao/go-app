package controller

import (
	"catalog/controller/boardgame"
	"catalog/controller/category"
	"catalog/controller/mechanism"
	"catalog/controller/tag"
	"catalog/repository"
)

// Controllers contains all the controllers
type Controllers struct {
	BoardgameController *boardgame.Controller
	TagController       *tag.Controller
	CategoryController  *category.Controller
	MechanismController *mechanism.Controller
}

// InitControllers returns a new Controllers
func InitControllers(repositories *repository.Repositories) *Controllers {
	return &Controllers{
		BoardgameController: boardgame.InitController(repositories.BoardGameRepository, repositories.TagRepository, repositories.CategoryRepository, repositories.MechanismRepository),
		TagController:       tag.InitController(repositories.TagRepository),
		CategoryController:  category.InitController(repositories.CategoryRepository),
		MechanismController: mechanism.InitController(repositories.MechanismRepository),
	}
}
