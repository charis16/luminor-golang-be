package utils

import "strings"

func GenerateSlug(str string) string {
	return strings.ReplaceAll(str, " ", "-")
}
