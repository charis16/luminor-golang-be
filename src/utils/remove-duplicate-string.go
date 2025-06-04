package utils

func RemoveDuplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, s := range slice {
		if _, exists := seen[s]; !exists && s != "" {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
