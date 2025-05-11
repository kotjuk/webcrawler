package crawler

import (
	"strings"
	"webcrawler/index"
	"webcrawler/logger"
	"webcrawler/monitor"

	"github.com/gocolly/colly"
)

func Crawl(startURL string, idx *index.Index, mon *monitor.Monitor) {
	c := colly.NewCollector(
		colly.MaxDepth(2),
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

	c.Visit(startURL)
}
