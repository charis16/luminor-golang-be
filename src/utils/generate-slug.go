package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug membersihkan string agar menjadi slug URL-friendly
func GenerateSlug(str string) string {
	// Lowercase
	slug := strings.ToLower(str)

	// Ganti karakter khusus yang umum
	replacements := map[string]string{
		"&":  "and",
		"@":  "at",
		"#":  "",
		"%":  "percent",
		"$":  "dollar",
		"(":  "",
		")":  "",
		"[":  "",
		"]":  "",
		"{":  "",
		"}":  "",
		"'":  "",
		"\"": "",
		".":  "",
		",":  "",
		"/":  "-",
		"\\": "-",
		":":  "",
		";":  "",
	}

	for old, newVal := range replacements {
		slug = strings.ReplaceAll(slug, old, newVal)
	}

	// Ganti spasi dan underscore dengan strip
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Hapus karakter non-alfanumerik kecuali strip (-)
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Hapus strip ganda
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// Hapus strip di awal dan akhir
	slug = strings.Trim(slug, "-")

	return slug
}
