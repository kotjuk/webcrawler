package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"webcrawler/crawler"
	"webcrawler/index"
	"webcrawler/logger"
	"webcrawler/monitor"
	"webcrawler/search"
)

func main() {
	logger.Init()
	idx := index.New()
	mon := monitor.New()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите поисковый запрос: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	links, err := search.GetLinks(query)
	if err != nil || len(links) == 0 {
		log.Fatalf("Не удалось найти ссылки по запросу: %s", query)
	}

	fmt.Println("Найдено ссылок:", len(links))
	for _, link := range links {
		fmt.Println("Начинаем обход:", link)
		crawler.Crawl(link, idx, mon)
	}

	results := idx.Search(query)
	fmt.Println("Результаты поиска по словам из страницы:")
	for _, url := range results {
		fmt.Println(url)
	}

	mon.Report()
}
