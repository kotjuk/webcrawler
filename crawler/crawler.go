package crawler

import (
	"fmt"
	"strings"
	"webcrawler/index"
	"webcrawler/logger"
	"webcrawler/monitor"

	"github.com/gocolly/colly"
)

type Options struct {
	MaxDepth int
}

func (o Options) normalized() (Options, error) {
	if o.MaxDepth == 0 {
		o.MaxDepth = 2
	}
	if o.MaxDepth < 1 || o.MaxDepth > 10 {
		return Options{}, fmt.Errorf("maxDepth must be between 1 and 10")
	}
	return o, nil
}

func Crawl(startURL string, idx *index.Index, mon *monitor.Monitor, opts Options) error {
	opts, err := opts.normalized()
	if err != nil {
		return err
	}
	c := colly.NewCollector(
		colly.MaxDepth(opts.MaxDepth),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		e.Request.Visit(link)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		text := e.Text
		words := strings.Fields(text)
		idx.Add(words, e.Request.URL.String())
		logger.Log("Проиндексирована страница: " + e.Request.URL.String())
		mon.IncrementPages()
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Log("Посещение: " + r.URL.String())
	})

	return c.Visit(startURL)
}
