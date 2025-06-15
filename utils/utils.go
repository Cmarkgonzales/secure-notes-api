package utils

import (
	"strings"
)

func ContainsIgnoreCase(text string, search string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(search))
}
