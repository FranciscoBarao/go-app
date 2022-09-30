package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAreCharactersUnique(t *testing.T) {

	str := ""
	result := AreCharactersUnique(str)
	assert.Equal(t, result, true, "#1 -> Result should be true")

	str = "bar"
	result = AreCharactersUnique(str)
	assert.Equal(t, result, true, "#2 -> Result should be true")

	str = "foo"
	result = AreCharactersUnique(str)
	assert.Equal(t, result, false, "#3 -> Result should be false")
}
