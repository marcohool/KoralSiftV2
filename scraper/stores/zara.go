package stores

import (
	"KoralSiftV2/browser"
	"KoralSiftV2/helpers"
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
}

func ScrapeAllProductsPage(browserCtx context.Context, baseUrl string) []string {
	fmt.Printf("Scraping Zara products page %s", baseUrl)

	page := 1

	var allHrefs []string

	for {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		hrefs := GetProductsFromPage(tabCtx, baseUrl, page)
		tabCancel()

		if len(hrefs) == 0 {
			fmt.Printf("No more products on page %d", page)
			break
		}

		fmt.Println("Found", len(hrefs), "products on page", page)

		page++
		allHrefs = append(allHrefs, hrefs...)

		fmt.Printf("Next page %d", page)
	}

	hrefsSet := helpers.CreateSet(allHrefs)

	fmt.Println("Total found product hrefs:", len(allHrefs))
	fmt.Println("Unique product hrefs:", len(hrefsSet))

	return allHrefs
}

func GetProductsFromPage(browserCtx context.Context, baseURL string, pageNo int) []string {
	url := fmt.Sprintf("%s&page=%d", baseURL, pageNo)
	fmt.Printf("\nScraping Zara page %d: %s\n", pageNo, url)

	ctx, cancel := context.WithTimeout(browserCtx, 20*time.Second)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML("html", &html),
	)
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
