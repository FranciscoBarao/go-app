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
	repo repository
}

// InitController initializes the boargame controller.
func InitController(boardGameRepo *boardgameRepo.BoardGameRepository) *Controller {
	return &Controller{
		repo: boardGameRepo,
	}
}

func (controller *Controller) Create(w http.ResponseWriter, r *http.Request) {

	temp := model.NewBoardGame("Bilbo", "Baggings", 10.0, 4)

	err := controller.repo.Create(temp)
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	render.New().JSON(w, http.StatusOK, "200")
}

func (controller *Controller) Get(w http.ResponseWriter, r *http.Request) {

	boardgame, err := controller.repo.GetByName("Bilbo")
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

func (controller *Controller) Update(w http.ResponseWriter, r *http.Request) {

	boardgame, err := controller.repo.GetByName("Bilbo")
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	boardgame.Dealer = "420"
	err = controller.repo.Update(boardgame)
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}
	render.New().JSON(w, http.StatusOK, "200 - Updated")
}

func (controller *Controller) Delete(w http.ResponseWriter, r *http.Request) {

	boardgame, err := controller.repo.GetByName("Bilbo")
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}

	err = controller.repo.DeleteById(boardgame)
	if err != nil {
		render.New().JSON(w, http.StatusOK, "500")
	}
	render.New().JSON(w, http.StatusOK, "200 - Deleted")
}
