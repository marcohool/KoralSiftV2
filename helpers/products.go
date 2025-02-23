package helpers

import (
	"KoralSiftV2/browser"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func ScrapeAllProductsPage(
	browserCtx context.Context,
	baseUrl string,
	hrefElementSelector string,
	urlPageFormatter func(string, int) string,
) []string {
	log.Info().Str("url", baseUrl).Msg("Scraping products page")

	page := 1
	hrefsMap := make(map[string]struct{})

	for {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		hrefs := getProductsFromPage(tabCtx, urlPageFormatter(baseUrl, page), hrefElementSelector)
		tabCancel()

		if len(hrefs) == 0 {
			log.Info().Int("page", page).Msg("No more products found")
			break
		}

		log.Info().Int("page", page).Int("count", len(hrefs)).Msg("Found products on page")

		for _, href := range hrefs {
			hrefsMap[href] = struct{}{}
		}

		page++
	}

	uniqueHrefs := CreateSliceFromMapKey(hrefsMap)

	log.Info().Int("unique_count", len(uniqueHrefs)).Msg("Total unique product hrefs found")

	return uniqueHrefs
}

func getProductsFromPage(browserCtx context.Context, url string, selector string) []string {
	log.Debug().Str("url", url).Msg("Scraping All Products page")

	ctx, cancel := context.WithTimeout(browserCtx, 20*time.Second)
	defer cancel()

	err, html := browser.ScrapePage(ctx, url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error navigating page")
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error parsing HTML")
		return nil
	}

	var productHrefs []string

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			productHrefs = append(productHrefs, href)
		} else {
			log.Warn().Str("url", url).Str("selector", selector).Msg("Failed to find href attribute")
		}
	})

	return productHrefs
}
