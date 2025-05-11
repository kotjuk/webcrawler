package index

import (
	"strings"
	"sync"
)

type Index struct {
	data map[string]map[string]struct{}
	mu   sync.RWMutex
}

func New() *Index {
	return &Index{
		data: make(map[string]map[string]struct{}),
	}
}

func (idx *Index) Add(words []string, url string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	for _, word := range words {
		word = strings.ToLower(word)
		if idx.data[word] == nil {
			idx.data[word] = make(map[string]struct{})
		}
		idx.data[word][url] = struct{}{}
	}
}

func (idx *Index) Search(query string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	urlsMap := idx.data[strings.ToLower(query)]
	var urls []string
	for url := range urlsMap {
		urls = append(urls, url)
	}
	return urls
}
