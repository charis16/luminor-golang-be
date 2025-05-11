package utils

import (
	"net/url"
	"path"
)

func GetObjectNameFromURL(fileURL string) string {
	u, err := url.Parse(fileURL)
	if err != nil {
		return fileURL // fallback
	}
	return path.Base(u.Path) // ambil hanya nama file
}
