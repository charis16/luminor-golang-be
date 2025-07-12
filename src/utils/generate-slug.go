package utils

import "strings"

func GenerateSlug(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
}
