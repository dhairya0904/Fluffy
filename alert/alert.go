package alert

import (
	"fluffy/domain"
	"fmt"
	"math"
	"sync"
	"time"
)

type alertMonitor struct {
	events  []domain.Event
	mu      sync.RWMutex
	onAlert bool
}

func NewAlertMonitor() *alertMonitor {
	events := make([]domain.Event, 0)
	return &alertMonitor{events: events}
}

func (monitor *alertMonitor) StartAlertMonitor(events chan domain.Event, done chan bool) {

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case event := <-events:
			monitor.consume(event)
		case <-ticker.C:
			go monitor.isOnAlert()
		case <-done:
			return
		}
	}
}

func (monitor *alertMonitor) consume(event domain.Event) {
	monitor.mu.Lock()
	monitor.events = append(monitor.events, event)
	monitor.mu.Unlock()
}

func (monitor *alertMonitor) isOnAlert() {

	if len(monitor.events) == 0 || monitor.events[0].Time.After(time.Now().UTC().Add(-5*time.Second)) {
		return
	}

	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	index := getIndexBeforeGivenTimeStamp(monitor.events, time.Now().UTC().Add(-5*time.Second))

	requestCount := len(monitor.events) - index
	requestPerSecond := math.Ceil(float64(requestCount) / float64(5))

	monitor.cleanUp(index)

	if requestPerSecond > 10 {
		monitor.onAlert = true
		fmt.Printf("Monitor on alert %f at time %s\n", requestPerSecond, time.Now().UTC().String())
	} else if monitor.onAlert {
		monitor.onAlert = false
		fmt.Printf("Monitor out of alert %f at %s\n", requestPerSecond, time.Now().UTC().String())
	}
}

func (monitor *alertMonitor) cleanUp(index int) {
	monitor.events = monitor.events[index:]
}

//// It will get the index of the first element that is before 2 minutes
func getIndexBeforeGivenTimeStamp(events []domain.Event, twoMinutes time.Time) int {

	fmt.Println(events)
	fmt.Println(twoMinutes)

	low, high := 0, len(events)-1
	index := -1

	for low <= high {

		mid := low + (high-low)/2

		if events[mid].Time.After(twoMinutes) {
			high = mid - 1
		} else {
			index = mid
			low = mid + 1
		}

	}
	fmt.Println(index)
	fmt.Println(events[index])
	return index
}
