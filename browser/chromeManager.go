package browser

import (
	"context"
	"github.com/chromedp/chromedp"
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
