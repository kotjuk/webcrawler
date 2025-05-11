package monitor

import (
	"fmt"
	"sync"
)

type Monitor struct {
	pagesCrawled int
	mu           sync.Mutex
}

func New() *Monitor {
	return &Monitor{}
}

func (m *Monitor) IncrementPages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pagesCrawled++
}

func (m *Monitor) Report() {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Printf("Обход завершен. Страниц обработано: %d\n", m.pagesCrawled)
}
