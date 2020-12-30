package monitor

import (
	"context"
	"fluffy/domain"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type Monitor struct {
	events              []domain.Event // temporary storage for incoming events
	mu                  sync.RWMutex   // Locks
	errorCount          int            // Metric to store failed requests
	out                 *os.File       // output file, where reports will be written
	reportTimeInSeconds time.Duration  // Time in seconds, report will be generated after every {reportTimeInSeconds}
}

func NewMonitor(reportTime int) *Monitor {

	events := make([]domain.Event, 0)
	return &Monitor{events: events, errorCount: 0, reportTimeInSeconds: time.Duration(reportTime) * time.Second}
}

func (monitor *Monitor) StartMonitor(ctx context.Context, wg *sync.WaitGroup, events chan domain.Event, reportPath string) {

	f, err := os.Create(reportPath)
	ticker := time.NewTicker(monitor.reportTimeInSeconds)

	if err != nil {
		log.Print("Error creating file for reports", err)
		os.Exit(1)
	}
	monitor.out = f

	defer monitor.out.Close()
	defer wg.Done()
	defer ticker.Stop()

	f.WriteString(fmt.Sprintf(DOG))

monitor:
	for {
		select {
		case event := <-events:
			monitor.consume(event)
		case <-ticker.C:
			go monitor.show()
		case <-ctx.Done():
			break monitor
		}
	}
}

func (monitor *Monitor) consume(event domain.Event) {

	monitor.mu.Lock()
	defer monitor.mu.Unlock()
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
	monitor.out.WriteString(fmt.Sprintf("\n---ERROR RATE: %.2f%%----\n\n\n", errorRate))
	// fmt.Printf("\n---ERROR RATE: %.2f%%----\n", errorRate)
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

	monitor.out.WriteString(fmt.Sprintf(REPORT_MAX_HITS, date))
	// fmt.Printf(REPORT_MAX_HITS, date)

	for i, hit := range urlCount {
		if i > 2 {
			break
		}
		monitor.out.WriteString(fmt.Sprintf(ROW, hit.Key, hit.Value))
		// fmt.Printf(ROW, hit.Key, hit.Value)
	}
}
