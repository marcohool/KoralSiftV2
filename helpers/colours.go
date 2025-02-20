package helpers

import (
	"KoralSiftV2/models"
	"fmt"
	"regexp"
	"strconv"
)

func ExtractHexFromStyle(style string) (string, bool) {
	reHex := regexp.MustCompile(`#([0-9A-Fa-f]{6})`)

	if match := reHex.FindString(style); match != "" {
		return match, true
	}

	return "", false
}

func RgbToHex(rgb string) string {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(rgb, -1)

	if len(matches) != 3 {
		return ""
	}

	r, _ := strconv.Atoi(matches[0])
	g, _ := strconv.Atoi(matches[1])
	b, _ := strconv.Atoi(matches[2])

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// MergeColours merges two slices of Colour structs.
// If all fields (Name, Hex, or ImageUrl) match, the structs are merged.
func MergeColours(slice1, slice2 []models.Colour) []models.Colour {
	merged := make([]models.Colour, 0)
	seen := make(map[string]models.Colour)

	for _, colour := range append(slice1, slice2...) {
		key := colour.Name + "|" + colour.Hex + "|" + colour.ImageUrl

		if existing, exists := seen[key]; !exists {
			seen[key] = existing
		}
	}

	return merged
}
