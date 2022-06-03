package controller

import (
	"go-app/controller/boardgame"
	"go-app/controller/category"
	"go-app/controller/tag"
	"go-app/repository"
)

// Controllers contains all the controllers
type Controllers struct {
	BoardgameController *boardgame.Controller
	TagController       *tag.Controller
	CategoryController  *category.Controller
}

// InitControllers returns a new Controllers
func InitControllers(repositories *repository.Repositories) *Controllers {
	return &Controllers{
		BoardgameController: boardgame.InitController(repositories.BoardGameRepository, repositories.TagRepository, repositories.CategoryRepository),
		TagController:       tag.InitController(repositories.TagRepository),
		CategoryController:  category.InitController(repositories.CategoryRepository),
	}
}
