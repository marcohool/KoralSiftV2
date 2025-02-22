package browser

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

type ChromeManager struct {
	allocCtx   context.Context
	browserCtx context.Context
	cancel     context.CancelFunc
}

func NewChromeManager() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/115.0.0.0 Safari/537.36"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	return browserCtx, func() { allocCancel(); browserCancel() }
}

func ScrapePage(browserCtx context.Context, url string) (error, string) {
	var html string

	err := chromedp.Run(browserCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML("html", &html),
	)

	return err, html
}
