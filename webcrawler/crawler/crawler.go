package crawler

import (
	"fmt"
	"net/http"
	"time"

	"webcrawler/parser"
)

// Стартовая функция краулера
func StartCrawling(query string) []string {
	// Начнём с одного стартового сайта
	startUrls := []string{
		"https://en.wikipedia.org/wiki/" + query, // примитивно
	}

	var foundLinks []string

	for _, url := range startUrls {
		fmt.Println("Crawling:", url)

		body, err := fetch(url)
		if err != nil {
			fmt.Println("Ошибка загрузки:", err)
			continue
		}

		links := parser.ExtractLinks(body, url)
		foundLinks = append(foundLinks, links...)
		time.Sleep(1 * time.Second) // не спамим
	}

	return foundLinks
}

func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("не 200 OK: %d", resp.StatusCode)
	}

	buf := make([]byte, resp.ContentLength)
	resp.Body.Read(buf)
	return string(buf), nil
}
