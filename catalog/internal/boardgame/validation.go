package boardgame

import (
	"net/http"

	"github.com/FranciscoBarao/catalog/internal/middleware"
)

func validatePlayersOrder(min, max int) error {
	if min > max {
		return middleware.NewError(http.StatusBadRequest, "min_players cannot exceed max_players")
	}
	return nil
}

func validatePlayTimeOrder(min, max int) error {
	if min == 0 || max == 0 {
		return nil
	}
	if min > max {
		return middleware.NewError(http.StatusBadRequest, "min_play_time cannot exceed max_play_time")
	}
	return nil
}
