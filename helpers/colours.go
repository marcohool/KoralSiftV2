package helpers

import (
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
