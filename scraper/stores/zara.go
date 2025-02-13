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

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/115.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
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
