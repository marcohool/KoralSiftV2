package scraper

import (
	"context"
	"github.com/chromedp/chromedp"
)

type ChromeManager struct {
	allocCtx   context.Context
	browserCtx context.Context
	cancel     context.CancelFunc
}

func NewChromeManager() *ChromeManager {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) "+
			"Chrome/115.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	return &ChromeManager{
		allocCtx:   allocCtx,
		browserCtx: browserCtx,
		cancel:     func() { cancel(); browserCancel() },
	}
}
