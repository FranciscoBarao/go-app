package slug

import (
	"fmt"
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// FromName derives a URL-safe slug from a display name.
func FromName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if name == "" {
		return ""
	}

	return s
}

// WithSuffix appends a numeric suffix for collision resolution.
func WithSuffix(base string, n int) string {
	if n <= 1 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, n)
}
