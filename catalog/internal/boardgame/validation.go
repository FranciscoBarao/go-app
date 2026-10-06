package boardgame

import "fmt"

func validatePlayersOrder(min, max int) error {
	if min > max {
		return fmt.Errorf("min_players cannot exceed max_players: %w", ErrValidation)
	}
	return nil
}

func validatePlayTimeOrder(min, max int) error {
	if min == 0 || max == 0 {
		return nil
	}
	if min > max {
		return fmt.Errorf("min_play_time cannot exceed max_play_time: %w", ErrValidation)
	}
	return nil
}
