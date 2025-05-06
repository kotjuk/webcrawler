package parser

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Извлекает все ссылки с HTML-страницы
func ExtractLinks(body string, base string) []string {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil
	}

	baseUrl, _ := url.Parse(base)
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					absolute := toAbsoluteURL(link, baseUrl)
					if absolute != "" {
						links = append(links, absolute)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}

func toAbsoluteURL(href string, base *url.URL) string {
	u, err := url.Parse(href)
	if err != nil || u.Scheme == "mailto" || u.Scheme == "javascript" {
		return ""
	}
	return base.ResolveReference(u).String()
}
