package exercises

/*
Logic:

Use a Map, key -> Character and value -> boolean
If we try to add a second equal character then it will already be true
*/
func AreCharactersUnique(str string) bool {

	characterMap := map[rune]bool{}

	for _, character := range str {

		if characterAlreadyExists(characterMap, character) {
			return false
		}

		characterMap[character] = true
	}

	return true
}

func characterAlreadyExists(mp map[rune]bool, character rune) bool {
	return mp[character]
}
