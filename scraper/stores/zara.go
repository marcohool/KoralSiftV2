package stores

import (
	"KoralSiftV2/browser"
	"KoralSiftV2/helpers"
	"KoralSiftV2/models"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

func ScrapeZara() {
	fmt.Println("Start Zara Scraper")

	browserCtx, cancel := browser.NewChromeManager()
	defer cancel()

	ukMenHrefs := ScrapeAllProductsPage(browserCtx, "https://www.zara.com/uk/en/man-all-products-l7465.html?v1=2443335")

	for _, href := range ukMenHrefs {
		fmt.Println(href)
	}

	ukMensProducts := ScrapeProductHrefs(browserCtx, ukMenHrefs, "Male", "GBP", "UK")

	for _, product := range ukMensProducts {
		fmt.Printf("%+v\n", product)
	}
}

func ScrapeProductHrefs(browserCtx context.Context,
	hrefs []string,
	gender string,
	currencyCode string,
	sourceRegion string) []*models.ClothingItem {
	fmt.Println("Scraping Zara product hrefs")

	var allProducts []*models.ClothingItem

	for _, href := range hrefs {
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
		}
	}

	return allProducts
}

func ScrapeProduct(browserCtx context.Context,
	url string) *models.ClothingItem {
	fmt.Println("Scraping Zara product at URL: ", url)

	err, html := browser.ScrapePage(browserCtx, url)
	if err != nil {
		fmt.Println("Error navigating page:", err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil
	}

	clothingItem := &models.ClothingItem{}

	clothingItem.Name = doc.Find(".product-detail-info__header-name").Text()

	price, err := helpers.ExtractNumericValue(doc.Find(".money-amount__main").Text())
	if err != nil {
		fmt.Println("Error extracting price:", err)
		return nil
	}
	clothingItem.Price = price

	imageURL, imgExists := doc.Find(".product-detail-view__main-image img").Attr("src")
	if !imgExists {
		fmt.Println("No image found for product at URL:", url)
	}
	clothingItem.ImageUrl = imageURL

	var hexColors []string
	doc.Find(".product-detail-color-selector__color-area").Each(func(i int, s *goquery.Selection) {
		colourStyle, colourExists := s.Attr("style")
		if colourExists {
			rgb, rgbExists := helpers.ExtractRGBFromStyle(colourStyle)
			if rgbExists {
				hexColor := helpers.RgbToHex(rgb)
				hexColors = append(hexColors, hexColor)
			}
		}
	})
	clothingItem.Colours = hexColors

	fmt.Println("Scraped product: ", clothingItem.Name)

	return clothingItem
}

func ScrapeAllProductsPage(browserCtx context.Context, baseUrl string) []string {
	fmt.Printf("Scraping Zara products page %s", baseUrl)

	page := 1

	hrefsMap := make(map[string]struct{})

	for {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		hrefs := GetProductsFromPage(tabCtx, baseUrl, page)
		tabCancel()

		if len(hrefs) == 0 {
			fmt.Printf("No more products on page %d", page)
			break
		}

		fmt.Println("Found", len(hrefs), "products on page", page)

		for _, href := range hrefs {
			hrefsMap[href] = struct{}{}
		}

		page++

		fmt.Printf("Next page %d", page)
	}

	uniqueHrefs := helpers.CreateSliceFromMap(hrefsMap)

	fmt.Println("Unique product hrefs:", len(uniqueHrefs))

	return uniqueHrefs
}

func GetProductsFromPage(browserCtx context.Context, baseURL string, pageNo int) []string {
	url := fmt.Sprintf("%s&page=%d", baseURL, pageNo)
	fmt.Printf("\nScraping Zara page %d: %s\n", pageNo, url)

	ctx, cancel := context.WithTimeout(browserCtx, 20*time.Second)
	defer cancel()

	err, html := browser.ScrapePage(ctx, url)
	if err != nil {
		fmt.Println("Error navigating page:", err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
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
