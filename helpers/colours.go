package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
ExtractColourFromStyle extracts a colour from a CSS style attribute.
It supports both hex and rgb colours, and returns the hex value.
*/
func ExtractColourFromStyle(style string) (string, bool) {
	hexRegex := regexp.MustCompile(`#([0-9a-fA-F]{6})`)
	rgbRegex := regexp.MustCompile(`rgb\((\d+),\s*(\d+),\s*(\d+)\)`)

	if hexMatch := hexRegex.FindString(style); hexMatch != "" {
		return strings.ToLower(hexMatch), true
	}

	if rgbMatch := rgbRegex.FindStringSubmatch(style); len(rgbMatch) == 4 {
		r, _ := strconv.Atoi(rgbMatch[1])
		g, _ := strconv.Atoi(rgbMatch[2])
		b, _ := strconv.Atoi(rgbMatch[3])

		hexColor := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		return hexColor, true
	}

	return "", false
}

//// MergeColours merges two slices of Colour structs.
//// If all fields (Name, Hex, or ImageUrl) match, the structs are merged.
//func MergeColours(slice1, slice2 []models.Colour) []models.Colour {
//	merged := make([]models.Colour, 0)
//	seen := make(map[string]models.Colour)
//
//	for _, colour := range append(slice1, slice2...) {
//		key := colour.Name + "|" + colour.Hex + "|" + colour.ImageUrl
//
//		if existing, exists := seen[key]; !exists {
//			merged = append(merged, colour)
//			seen[key] = existing
//		}
//	}
//
//	return merged
//}
