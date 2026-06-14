package ai

import (
	"regexp"
	"strconv"
	"strings"
)

var bodyCountPattern = regexp.MustCompile(`(\d+)\s*[-]?\s*cuerp`)

var spanishBodyPhrases = []struct {
	phrase string
	count  int
}{
	{"seis cuerpos", 6},
	{"seis cuerpo", 6},
	{"cinco cuerpos", 5},
	{"cinco cuerpo", 5},
	{"cuatro cuerpos", 4},
	{"cuatro cuerpo", 4},
	{"tres cuerpos", 3},
	{"tres cuerpo", 3},
	{"dos cuerpos", 2},
	{"dos cuerpo", 2},
	{"un cuerpo", 1},
	{"uno cuerpo", 1},
	{"una cuerpo", 1},
}

func parseBodyCount(description string) int {
	lower := strings.ToLower(description)

	if match := bodyCountPattern.FindStringSubmatch(lower); len(match) >= 2 {
		if count, err := strconv.Atoi(match[1]); err == nil {
			return clampBodyCount(count)
		}
	}

	for _, entry := range spanishBodyPhrases {
		if strings.Contains(lower, entry.phrase) {
			return entry.count
		}
	}

	return 2
}

func clampBodyCount(count int) int {
	if count < 1 {
		return 1
	}
	if count > 6 {
		return 6
	}
	return count
}

func equalRatios(count int) []float64 {
	ratios := make([]float64, count)
	share := 1.0 / float64(count)
	for i := range ratios {
		ratios[i] = share
	}
	return ratios
}

type deskDrawerIntent struct {
	enabled  bool
	mode     string // "single" | "tower"
	count    int
	bayRatio float64
	legRatio float64
}

var drawerCountPattern = regexp.MustCompile(`(\d+)\s*[-]?\s*cajon`)

var spanishDrawerPhrases = []struct {
	phrase string
	count  int
}{
	{"tres cajones", 3},
	{"tres cajon", 3},
	{"dos cajones", 2},
	{"dos cajon", 2},
	{"un cajon", 1},
	{"un cajón", 1},
	{"una cajon", 1},
	{"1 cajon", 1},
	{"1 cajón", 1},
}

func parseDrawerCount(description string) int {
	lower := strings.ToLower(description)

	if match := drawerCountPattern.FindStringSubmatch(lower); len(match) >= 2 {
		if count, err := strconv.Atoi(match[1]); err == nil {
			return clampDrawerCount(count)
		}
	}

	for _, entry := range spanishDrawerPhrases {
		if strings.Contains(lower, entry.phrase) {
			return entry.count
		}
	}

	return 0
}

func clampDrawerCount(count int) int {
	if count < 1 {
		return 1
	}
	if count > 6 {
		return 6
	}
	return count
}

func hasDrawerKeyword(description string) bool {
	lower := strings.ToLower(description)
	if strings.Contains(lower, "sin cajon") {
		return false
	}
	if strings.Contains(lower, "cajonera") {
		return true
	}
	if parseDrawerCount(description) > 0 {
		return true
	}
	for _, entry := range spanishDrawerPhrases {
		if strings.Contains(lower, entry.phrase) {
			return true
		}
	}
	return false
}

func parseDeskDrawerIntent(description string) deskDrawerIntent {
	if !hasDrawerKeyword(description) {
		return deskDrawerIntent{}
	}

	lower := strings.ToLower(description)
	count := parseDrawerCount(description)

	if strings.Contains(lower, "cajonera") || count > 1 {
		if count < 2 {
			count = 3
		}
		return deskDrawerIntent{
			enabled:  true,
			mode:     "tower",
			count:    count,
			bayRatio: 0.28,
			legRatio: 0.72,
		}
	}

	return deskDrawerIntent{
		enabled:  true,
		mode:     "single",
		count:    1,
		bayRatio: 0.18,
		legRatio: 0.82,
	}
}
