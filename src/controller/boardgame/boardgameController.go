package boardgame

import (
	"net/http"
	"strconv"

	"github.com/unrolled/render"

	"go-app/model"
	"go-app/repository/boardgameRepo"
)

// declaring the repository interface in the controller package allows us to easily swap out the actual implementation, enforcing loose coupling.
type repository interface {
	Create(boardgame model.BoardGame) error
	GetByName(name string) (model.BoardGame, error)
	Update(boardgame model.BoardGame) error
	DeleteById(boardgame model.BoardGame) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type Controller struct {
	service repository
}

// InitController initializes the boargame controller.
func InitController(boardGameRepo *boardgameRepo.BoardGameRepository) *Controller {
	return &Controller{
		service: boardGameRepo,
	}
}

func (controller *Controller) CreateBG(w http.ResponseWriter, r *http.Request) {

	temp := model.NewBoardGame("Bilbo", "Baggings", 10.0, 4)

	err := controller.service.Create(temp)
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	render.New().JSON(w, http.StatusOK, "200")
}

func (controller *Controller) GetBG(w http.ResponseWriter, r *http.Request) {

	boardgame, err := controller.service.GetByName("Bilbo")
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	response := map[string]string{
		"Name":         boardgame.GetName(),
		"Dealer":       boardgame.GetDealer(),
		"Price":        strconv.FormatFloat(boardgame.GetPrice(), 'g', 1, 64),
		"PlayerNumber": strconv.Itoa(boardgame.GetPlayerNumber()),
	}

	render.New().JSON(w, http.StatusOK, response)
}

func (controller *Controller) UpdateBG(w http.ResponseWriter, r *http.Request) {

	boardgame, err := controller.service.GetByName("Bilbo")
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	boardgame.Dealer = "420"
	err = controller.service.Update(boardgame)
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}
	render.New().JSON(w, http.StatusOK, "200 - Updated")
}

func (controller *Controller) DeleteBG(w http.ResponseWriter, r *http.Request) {

	boardgame, err := controller.service.GetByName("Bilbo")
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	err = controller.service.DeleteById(boardgame)
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}
	render.New().JSON(w, http.StatusOK, "200 - Deleted")
}
