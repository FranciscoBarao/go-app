package utils

import "regexp"

func StringInSlice(value string, list []string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// Method that checks if a string is alphanumeric
func IsAlphanumeric(word string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(word)
}
