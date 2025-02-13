package stores

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

func ScrapeZara() {
	fmt.Println("Start Zara Scraper")

	ScrapeProductsPage("https://www.zara.com/uk/en/man-all-products-l7465.html?v1=2443335")
}

func ScrapeProductsPage(baseUrl string) {
	fmt.Printf("Scraping Zara products page %s", baseUrl)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/115.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	defer browserCancel()

	page := 1

	for {
		tabCtx, tabCancel := chromedp.NewContext(browserCtx)
		hrefs := GetProductsFromPage(tabCtx, baseUrl, page)
		tabCancel()

		if len(hrefs) == 0 {
			fmt.Printf("No more products on page %d", page)
			break
		}

		fmt.Println("page scraped")
		page++

		fmt.Printf("Next page %d", page)
	}

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
			fmt.Printf("Product %d href: %s\n", i, href)
			productHrefs = append(productHrefs, href)
		}
	})

	return productHrefs
}
