package scraper

import (
	"KoralSiftV2/scraper/stores"
	"github.com/rs/zerolog/log"
)

func RunScraper() {
	log.Info().Msg("Running scrapers")

	//stores.ScrapeZara()
	stores.ScrapeAsos()
}
