package scraper

import (
	"KoralSiftV2/scraper/stores/asos"
	"KoralSiftV2/scraper/stores/zara"
	"github.com/rs/zerolog/log"
)

func RunScraper() {
	log.Info().Msg("Running scrapers")

	zara.ScrapeZara()
	asos.ScrapeAsos()
}
