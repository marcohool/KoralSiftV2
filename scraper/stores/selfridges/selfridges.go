package selfridges

import (
	"KoralSiftV2/browser"
	"KoralSiftV2/helpers"
	"fmt"
	"github.com/rs/zerolog/log"
)

func ScrapeSelfridges() {
	log.Info().Msg("Starting Selfridges Scraper")

	browserCtx, cancel := browser.NewChromeManager()
	defer cancel()

	// Do shoes as well -> https://www.selfridges.com/GB/en/cat/shoes/mens/
	ukMenHrefs := helpers.ScrapeAllProductsPage(
		browserCtx,
		"https://www.selfridges.com/GB/en/cat/mens/clothing/",
		"[data-analytics-link=\"product_card_link\"]",
		formatProductPageUrl,
	)

	for i, href := range ukMenHrefs {
		log.Info().Str("href", href).Msgf("Found product %d", i)
	}

}

func formatProductPageUrl(baseUrl string, page int) string {
	return fmt.Sprintf("%s?pn=%d", baseUrl, page)
}
