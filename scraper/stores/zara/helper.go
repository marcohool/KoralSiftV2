package zara

import (
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"strings"
)

func ExtractVariantInfo(doc *goquery.Document, url string) models.ProductVariant {
	imageURL, imgExists := extractImageURL(doc)
	if !imgExists {
		log.Warn().Str("url", url).Msg("No image found for product")
	}

	price, err := helpers.ExtractNumericValue(doc.Find(".money-amount__main").Text())
	if err != nil {
		log.Warn().Err(err).Str("url", url).Msg("Error extracting price")
	}

	name := strings.Split(doc.Find(".product-color-extended-name").Text(), " | ")[0]
	if name == "" {
		log.Warn().Err(err).Str("url", url).Msg("Error extracting product name")
	}

	var hex string
	selectedColour := doc.Find(
		".product-detail-color-selector__color-button--is-selected .product-detail-color-selector__color-area",
	)
	if selectedColour.Length() > 0 {
		colourStyle, _ := selectedColour.Attr("style")

		var hexExists bool
		hex, hexExists = helpers.ExtractColourFromStyle(colourStyle)
		if !hexExists {
			log.Warn().
				Str("style", colourStyle).
				Str("url", url).
				Msg("No hex colour found from style for product")
		}
	} else {
		log.Warn().Str("url", url).Msg("No selected colour found")
	}

	return models.NewProductVariant(name, hex, price, imageURL, url)
}

func extractImageURL(doc *goquery.Document) (string, bool) {
	return doc.Find("[class^='product-detail-view'] [class$='-image'] img").
		Attr("src")
}

func DeclineCookies(browserCtx context.Context) error {
	return chromedp.Run(browserCtx,
		chromedp.WaitVisible("#onetrust-reject-all-handler", chromedp.ByID),
		chromedp.ScrollIntoView("#onetrust-reject-all-handler"),
		chromedp.Click("#onetrust-reject-all-handler", chromedp.ByID),
		chromedp.WaitNotVisible("#onetrust-reject-all-handler", chromedp.ByID),
	)
}
