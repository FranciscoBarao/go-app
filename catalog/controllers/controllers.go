package controllers

import (
	"github.com/FranciscoBarao/catalog/services"
)

// Controllers contains all the controllers
type Controllers struct {
	BoardgameController *BoardgameController
	TagController       *TagController
	CategoryController  *CategoryController
	MechanismController *MechanismController
}

// InitControllers returns a new Controllers
func InitControllers(services *services.Services) *Controllers {
	return &Controllers{
		BoardgameController: InitBoardgameController(services.BoardgameService),
		TagController:       InitTagController(services.TagService),
		CategoryController:  InitCategoryController(services.CategoryService),
		MechanismController: InitMechanismController(services.MechanismService),
	}
}
