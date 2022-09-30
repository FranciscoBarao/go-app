package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArePermutations(t *testing.T) {

	str1, str2 := "", "a"
	result := ArePermutations(str1, str2)
	assert.Equal(t, result, false, "#1 -> Result should be false")

	str1, str2 = "a", ""
	result = ArePermutations(str1, str2)
	assert.Equal(t, result, false, "#2 -> Result should be false")

	str1, str2 = "Nib", "bin"
	result = ArePermutations(str1, str2)
	assert.Equal(t, result, false, "#3 -> Result should be false")

	str1, str2 = "act", "cat"
	result = ArePermutations(str1, str2)
	assert.Equal(t, result, true, "#4 -> Result should be true")

	str1, str2 = "a ct", "ca t"
	result = ArePermutations(str1, str2)
	assert.Equal(t, result, true, "#5 -> Result should be true")
}
