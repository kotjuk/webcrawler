package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"webcrawler/crawler"
	"webcrawler/index"
	"webcrawler/logger"
	"webcrawler/monitor"
	"webcrawler/search"
)

func main() {
	webMode := flag.Bool("web", false, "run web UI on http://localhost:8080")
	addr := flag.String("addr", "localhost:8080", "web UI listen address")
	flag.Parse()

	logger.Init()
	if *webMode {
		runWeb(*addr)
		return
	}

	runCLI()
}

func runCLI() {
	idx := index.New()
	mon := monitor.New()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите поисковый запрос: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	fmt.Print("Введите глубину обхода (1-10, по умолчанию 2): ")
	depthLine, _ := reader.ReadString('\n')
	depthLine = strings.TrimSpace(depthLine)
	depth := 2
	if depthLine != "" {
		if parsed, err := strconv.Atoi(depthLine); err == nil {
			depth = parsed
		}
	}

	links, err := search.GetLinks(query)
	if err != nil || len(links) == 0 {
		log.Fatalf("Не удалось найти ссылки по запросу: %s", query)
	}

	fmt.Println("Найдено ссылок:", len(links))
	for _, link := range links {
		fmt.Println("Начинаем обход:", link)
		if err := crawler.Crawl(link, idx, mon, crawler.Options{MaxDepth: depth}); err != nil {
			log.Printf("Ошибка обхода %s: %v", link, err)
		}
	}

	results := searchInIndex(idx, query)
	fmt.Println("Результаты поиска по словам из страниц:")
	for _, url := range results {
		fmt.Println(url)
	}

	mon.Report()
}

func searchInIndex(idx *index.Index, query string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, w := range strings.Fields(query) {
		for _, u := range idx.Search(w) {
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			out = append(out, u)
		}
	}
	return out
}

func runWeb(addr string) {
	tpl := template.Must(template.New("page").Parse(`<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>Webcrawler</title>
  <style>
    body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif; margin: 24px; color: #111; }
    .card { max-width: 900px; border: 1px solid #e5e5e5; border-radius: 12px; padding: 16px; }
    label { display: block; font-weight: 600; margin-top: 12px; }
    input { width: 100%; padding: 10px 12px; border-radius: 10px; border: 1px solid #ccc; font-size: 16px; }
    button { margin-top: 16px; padding: 10px 14px; border-radius: 10px; border: 0; background: #111; color: #fff; font-size: 16px; cursor: pointer; }
    code { background: #f6f6f6; padding: 2px 6px; border-radius: 6px; }
    .muted { color: #666; }
    ul { padding-left: 18px; }
    .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
    @media (max-width: 740px) { .grid { grid-template-columns: 1fr; } }
  </style>
</head>
<body>
  <div class="card">
    <h2 style="margin:0 0 6px 0;">Webcrawler</h2>
    <div class="muted">Локальный интерфейс. Логи пишутся в <code>crawler.log</code>.</div>

    <form method="post" action="/run">
      <label>Поисковый запрос</label>
      <input name="query" placeholder="например: golang web crawler" value="{{.Query}}" required />

      <div class="grid">
        <div>
          <label>Глубина обхода (1-10)</label>
          <input name="depth" inputmode="numeric" pattern="[0-9]*" value="{{.Depth}}" />
        </div>
        <div>
          <label>Адрес интерфейса</label>
          <input value="http://` + addr + `" readonly />
        </div>
      </div>

      <button type="submit">Запустить</button>
    </form>

    {{if .Error}}
      <p style="color:#b00020; font-weight:600;">Ошибка: {{.Error}}</p>
    {{end}}

    {{if .Links}}
      <h3>Seed‑ссылки ({{len .Links}})</h3>
      <ul>
        {{range .Links}}<li><a href="{{.}}" target="_blank" rel="noreferrer">{{.}}</a></li>{{end}}
      </ul>
    {{end}}

    {{if .Results}}
      <h3>Результаты ({{len .Results}})</h3>
      <ul>
        {{range .Results}}<li><a href="{{.}}" target="_blank" rel="noreferrer">{{.}}</a></li>{{end}}
      </ul>
    {{end}}

    {{if .PagesCrawled}}
      <p class="muted">Страниц обработано: <b>{{.PagesCrawled}}</b></p>
    {{end}}
  </div>
</body>
</html>`))

	type pageData struct {
		Query       string
		Depth       int
		Links       []string
		Results     []string
		PagesCrawled int
		Error       string
	}

	render := func(w http.ResponseWriter, d pageData) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tpl.Execute(w, d)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, pageData{Depth: 2})
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err := r.ParseForm(); err != nil {
			render(w, pageData{Depth: 2, Error: err.Error()})
			return
		}

		query := strings.TrimSpace(r.FormValue("query"))
		depthStr := strings.TrimSpace(r.FormValue("depth"))
		depth := 2
		if depthStr != "" {
			if v, err := strconv.Atoi(depthStr); err == nil {
				depth = v
			}
		}

		if query == "" {
			render(w, pageData{Query: query, Depth: depth, Error: "пустой запрос"})
			return
		}

		links, err := search.GetLinks(query)
		if err != nil || len(links) == 0 {
			msg := "не удалось найти ссылки по запросу"
			if err != nil {
				msg = err.Error()
			}
			render(w, pageData{Query: query, Depth: depth, Error: msg})
			return
		}

		idx := index.New()
		mon := monitor.New()
		for _, link := range links {
			if err := crawler.Crawl(link, idx, mon, crawler.Options{MaxDepth: depth}); err != nil {
				log.Printf("Ошибка обхода %s: %v", link, err)
			}
		}

		render(w, pageData{
			Query:        query,
			Depth:        depth,
			Links:        links,
			Results:      searchInIndex(idx, query),
			PagesCrawled: mon.Pages(),
		})
	})

	log.Printf("Web UI listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
