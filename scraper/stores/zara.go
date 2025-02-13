package stores

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"time"
)

func ScrapeZara() {
	fmt.Println("Start Zara Scraper")

	ScrapeProductsPage("https://www.zara.com/uk/en/man-all-products-l7465.html?v1=2443335")

}

func ScrapeProductsPage(url string) {
	fmt.Printf("Scraping Zara products page %s", url)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)

			if err != nil {
				return err
			}

			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)

			return err
		}),
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(html)
}
