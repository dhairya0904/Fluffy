package monitor

import (
	"fluffy/domain"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Monitor struct {
	events     []domain.Event
	mu         sync.RWMutex
	errorCount int
}

func NewMonitor() *Monitor {
	events := make([]domain.Event, 0)
	return &Monitor{events: events, errorCount: 0}
}

func (monitor *Monitor) StartMonitor(events chan domain.Event, done chan bool) {

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case event := <-events:
			monitor.mu.Lock()
			monitor.consume(event)
			monitor.mu.Unlock()
		case <-ticker.C:
			go monitor.show()
		case <-done:
			return
		}
	}
}

func (monitor *Monitor) consume(event domain.Event) {

	errorCount := event.Status / 100

	if errorCount == 5 {
		monitor.errorCount++
	}
	monitor.events = append(monitor.events, event)
}

func (monitor *Monitor) show() {

	if len(monitor.events) == 0 {
		return
	}

	monitor.mu.Lock()
	var oldEvents = make([]domain.Event, len(monitor.events))
	errorCount := monitor.errorCount
	copy(oldEvents, monitor.events)
	monitor.reset()
	monitor.mu.Unlock()
	monitor.showMostAccessedResource(oldEvents)
	monitor.showErrorRate(len(oldEvents), errorCount)
}

func (monitor *Monitor) reset() {
	monitor.events = nil
	monitor.errorCount = 0
}

func (monitor *Monitor) showErrorRate(events, errorCount int) {
	var errorRate float32 = float32(errorCount) / float32(events) * 100
	fmt.Printf("\n---ERROR RATE: %.2f%%----\n", errorRate)
}

func (monitor *Monitor) showMostAccessedResource(events []domain.Event) {

	count := make(map[string]int)
	date := events[0].Time.Format("2006-01-02 15:04:05")

	for _, event := range events {
		url := event.URL
		if val, ok := count[url]; ok {
			count[url] = val + 1
		} else {
			count[url] = 1
		}
	}

	type URLCount struct {
		Key   string
		Value int
	}

	var urlCount []URLCount
	for k, v := range count {
		urlCount = append(urlCount, URLCount{k, v})
	}

	sort.Slice(urlCount, func(i, j int) bool {
		return urlCount[i].Value > urlCount[j].Value
	})
	fmt.Println(DOG)
	fmt.Printf(REPORT_MAX_HITS, date)

	for i, hit := range urlCount {
		if i > 2 {
			break
		}
		fmt.Printf(ROW, hit.Key, hit.Value)
	}
}
