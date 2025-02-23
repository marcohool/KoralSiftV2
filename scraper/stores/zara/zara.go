package zara

import (
	"KoralSiftV2/browser"
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
	"KoralSiftV2/models/enums"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func ScrapeZara() {
	log.Info().Msg("Starting Zara Scraper")

	browserCtx, cancel := browser.NewChromeManager()
	defer cancel()

	ukMenHrefs := ScrapeZaraAllProductsPage(
		browserCtx,
		"https://www.zara.com/uk/en/man-all-products-l7465.html?v1=2443335",
	)

	ukMensProducts := ScrapeProductHrefs(browserCtx, ukMenHrefs, enums.Male, enums.GBP, enums.UK)

	err := helpers.SaveSliceToJSONFile(ukMensProducts, "zara")
	if err != nil {
		log.Error().Err(err).Msg("Failed to save Zara data")
	}
}

func ScrapeProductHrefs(browserCtx context.Context,
	hrefs []string,
	gender enums.Gender,
	currencyCode enums.CurrencyCode,
	sourceRegion enums.SourceRegion) []*models.ClothingItem {
	log.Debug().Int("count", len(hrefs)).Msg("Scraping Zara product hrefs")
	var allProducts []*models.ClothingItem

	for i, href := range hrefs {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		clothingItem := ScrapeProduct(tabCtx, href, currencyCode, gender, sourceRegion)
		tabCancel()

		if clothingItem != nil {
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
	url string,
	currencyCode enums.CurrencyCode,
	gender enums.Gender,
	region enums.SourceRegion) *models.ClothingItem {

	log.Debug().Str("url", url).Msg("Scraping Zara product")

	err, html := browser.ScrapePage(browserCtx, url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error navigating page")
		return nil
	}

	err = DeclineCookies(browserCtx)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error declining cookies")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error parsing HTML")
		return nil
	}

	name := doc.Find("h1[class^='product-detail-'][class$='name']").Text()
	if name == "" {
		log.Warn().Str("url", url).Msg("No product name found")
	}

	var productVariants []models.ProductVariant

	var colorNodes []*cdp.Node
	err = chromedp.Run(
		browserCtx,
		chromedp.Nodes(
			".product-detail-color-selector__color-button",
			&colorNodes,
			chromedp.AtLeast(0),
		),
	)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Error finding color nodes")
	}

	if len(colorNodes) == 0 {
		log.Debug().
			Str("product", name).
			Str("url", url).
			Msg("No colour selectors found for product on page")

		productVariants = append(productVariants, ExtractVariantInfo(doc, url))
	} else {
		for i := range colorNodes {
			var variantHtml string
			err := chromedp.Run(browserCtx,
				chromedp.MouseClickNode(colorNodes[i]),
				chromedp.Sleep(2*time.Second),
				chromedp.OuterHTML("html", &variantHtml),
			)
			if err != nil {
				log.Error().Err(err).Str("url", url).Msg("Failed to click color element")
				continue
			}

			variantDoc, err := goquery.NewDocumentFromReader(strings.NewReader(variantHtml))
			productVariants = append(productVariants, ExtractVariantInfo(variantDoc, url))
		}
	}

	clothingItem := models.NewClothingItem(
		name,
		nil,
		enums.Zara,
		productVariants,
		currencyCode,
		gender,
		region,
	)

	return &clothingItem
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

func ScrapeZaraAllProductsPage(browserCtx context.Context, baseUrl string) []string {
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

	uniqueHrefs := helpers.CreateSliceFromMapKey(hrefsMap)

	log.Info().Int("unique_count", len(uniqueHrefs)).Msg("Total unique product hrefs found")

	return uniqueHrefs
}
