package boardgame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreateBoardgameDTOUsesProvidedSlug(t *testing.T) {
	req := &CreateBoardgameRequest{
		Name:       "Settlers of Catan",
		MinPlayers: 2,
		MaxPlayers: 4,
	}

	input := newCreateBoardgameDTO(req, "settlers-of-catan", nil, nil, nil, nil)

	assert.Equal(t, "settlers-of-catan", input.Slug)
	assert.Equal(t, "Settlers of Catan", input.Name)
	assert.Equal(t, 2, input.MinPlayers)
}
