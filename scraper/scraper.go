package scraper

import (
	"KoralSiftV2/scraper/stores"
	"fmt"
)

func RunScraper() {
	fmt.Println("Running scrapers")

	// ScrapeZara()
	stores.ScrapeZara()
}
