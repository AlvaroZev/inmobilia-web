package ai

import "strings"

func isCloset(description string) bool {
	lower := strings.ToLower(description)
	return strings.Contains(lower, "closet") ||
		strings.Contains(lower, "ropero") ||
		strings.Contains(lower, "armario")
}

func isDesk(description string) bool {
	lower := strings.ToLower(description)
	return strings.Contains(lower, "escritorio") ||
		strings.Contains(lower, "desk") ||
		strings.Contains(lower, "mesa de trabajo") ||
		strings.Contains(lower, "home office")
}

func isEntertainmentCenter(description string) bool {
	lower := strings.ToLower(description)
	return strings.Contains(lower, "centro de entretenimiento") ||
		strings.Contains(lower, "mueble para tv") ||
		strings.Contains(lower, "mueble tv") ||
		strings.Contains(lower, "rack tv") ||
		strings.Contains(lower, "entertainment center") ||
		strings.Contains(lower, "tv stand") ||
		strings.Contains(lower, "modular tv")
}
