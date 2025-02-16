package helpers

import (
	"regexp"
	"strconv"
)

func ExtractNumericValue(priceText string) (float64, error) {
	re := regexp.MustCompile(`\d+\.\d+`)
	match := re.FindString(priceText)

	return strconv.ParseFloat(match, 64)
}
