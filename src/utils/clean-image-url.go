package utils

import "strings"

func CleanImageURL(url string) string {
	url = strings.Trim(url, `"`)
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}
