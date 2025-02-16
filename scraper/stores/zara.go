package stores

import (
	"KoralSiftV2/browser"
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func ScrapeZara() {
	log.Info().Msg("Starting Zara Scraper")

	browserCtx, cancel := browser.NewChromeManager()
	defer cancel()

	ukMenHrefs := ScrapeAllProductsPage(browserCtx, "https://www.zara.com/uk/en/man-all-products-l7465.html?v1=2443335")

	for _, href := range ukMenHrefs {
		log.Debug().Str("href", href).Msg("Found product href")
	}

	ukMensProducts := ScrapeProductHrefs(browserCtx, ukMenHrefs, "Male", "GBP", "UK")

	for _, product := range ukMensProducts {
		log.Info().Interface("product", product).Msg("Scraped product")
	}
}

func ScrapeProductHrefs(browserCtx context.Context,
	hrefs []string,
	gender string,
	currencyCode string,
	sourceRegion string) []*models.ClothingItem {
	log.Info().Int("count", len(hrefs)).Msg("Scraping Zara product hrefs")

	var allProducts []*models.ClothingItem

	for i, href := range hrefs {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		clothingItem := ScrapeProduct(tabCtx, href)
		tabCancel()

		if clothingItem != nil {
			clothingItem.Brand = "Zara"
			clothingItem.CurrencyCode = currencyCode
			clothingItem.Gender = gender
			clothingItem.SourceUrl = href
			clothingItem.SourceRegion = sourceRegion

			allProducts = append(allProducts, clothingItem)

			log.Info().
				Int("index", i+1).
				Int("total", len(hrefs)).
				Str("href", href).
				Interface("product", clothingItem).
				Msg("Scraped product successfully")
		} else {
			log.Warn().
				Int("index", i+1).
				Int("total", len(hrefs)).
				Str("href", href).
				Msg("Failed to scrape product")
		}
	}

	return allProducts
}

func ScrapeProduct(browserCtx context.Context,
	url string) *models.ClothingItem {
	log.Info().Str("url", url).Msg("Scraping Zara product")

	err, html := browser.ScrapePage(browserCtx, url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error navigating page")
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error parsing HTML")
		return nil
	}

	clothingItem := &models.ClothingItem{}

	clothingItem.Name = doc.Find(".product-detail-info__header-name").Text()

	price, err := helpers.ExtractNumericValue(doc.Find(".money-amount__main").Text())
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error extracting price")
		return nil
	}
	clothingItem.Price = price

	imageURL, imgExists := doc.Find(".product-detail-view__main-image img").Attr("src")
	if !imgExists {
		log.Warn().Str("url", url).Msg("No image found for product")
	}
	clothingItem.ImageUrl = imageURL

	hexColors := make([]string, 0)
	colorElements := doc.Find(".product-detail-color-selector__color-area")
	if colorElements.Length() == 0 {
		log.Warn().Str("product", clothingItem.Name).Str("url", url).Msg("No colour selectors found for product on page")
	}

	colorElements.Each(func(i int, s *goquery.Selection) {
		colourStyle, colourExists := s.Attr("style")
		if colourExists {
			hex, hexExists := helpers.ExtractHexFromStyle(colourStyle)
			if hexExists {
				hexColors = append(hexColors, hex)
			} else {
				log.Warn().Str("style", colourStyle).Str("url", url).Msg("No hex colour found from style for product")
			}
		} else {
			log.Warn().Str("url", url).Msg("No colours found for product")
		}
	})
	clothingItem.Colours = hexColors

	return clothingItem
}

func ScrapeAllProductsPage(browserCtx context.Context, baseUrl string) []string {
	log.Info().Str("url", baseUrl).Msg("Scraping Zara products page")

	page := 1

	hrefsMap := make(map[string]struct{})

	for {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		hrefs := GetProductsFromPage(tabCtx, baseUrl, page)
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

	uniqueHrefs := helpers.CreateSliceFromMap(hrefsMap)

	log.Info().Int("unique_count", len(uniqueHrefs)).Msg("Total unique product hrefs found")

	return uniqueHrefs
}

func GetProductsFromPage(browserCtx context.Context, baseURL string, pageNo int) []string {
	url := fmt.Sprintf("%s&page=%d", baseURL, pageNo)
	log.Debug().Int("page", pageNo).Str("url", url).Msg("Scraping Zara product page")

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

	doc.Find("li.product-grid-product a.product-link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			productHrefs = append(productHrefs, href)
		}
	})

	return productHrefs
}
