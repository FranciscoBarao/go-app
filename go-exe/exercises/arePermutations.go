package exercises

import "reflect"

/*
Logic:

Use a Map for each string, key -> Character and value -> int
If both maps are equal in the end, then they have the same characters and are therefore permutations
*/
func ArePermutations(str1, str2 string) bool {

	if str1 == "" || str2 == "" {
		return false
	}

	map1 := getCharacterMap(str1)
	map2 := getCharacterMap(str2)

	return reflect.DeepEqual(map1, map2)
}

func getCharacterMap(str string) map[rune]int {

	characterMap := map[rune]int{}
	for _, character := range str {
		characterMap[character] += 1
	}
	return characterMap
}
