package alert

import (
	"context"
	"fluffy/domain"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type alertMonitor struct {
	events      []domain.Event /// Temporary storage for storing events
	mu          sync.RWMutex   /// Mutex
	onAlert     bool           /// Flag to show if service is on alert
	out         *os.File       /// Output file where alerts will be written
	threshold   int            /// threshold(requests/sec) after which alarm will go off
	alertWindow time.Duration  /// Window to consider for calculating requests/seconds
}

func NewAlertMonitor(alertWindowInSeconds int, threshold int) *alertMonitor {
	events := make([]domain.Event, 0)
	return &alertMonitor{events: events, threshold: threshold, alertWindow: time.Duration(alertWindowInSeconds) * time.Second}
}

func (monitor *alertMonitor) StartAlertMonitor(ctx context.Context, wg *sync.WaitGroup, events chan domain.Event, alertPath string) {

	f, err := os.Create(alertPath)
	ticker := time.NewTicker(1 * time.Second)

	if err != nil {
		log.Print("Error creating file for alerts", err)
		os.Exit(1)
	}
	monitor.out = f
	defer monitor.out.Close()
	defer wg.Done()
	defer ticker.Stop()

alert:
	for {
		select {
		case event := <-events:
			monitor.consume(event)
		case <-ticker.C:
			go monitor.isOnAlert()
		case <-ctx.Done():
			break alert
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
	index := getIndexBeforeGivenTimeStamp(monitor.events, time.Now().UTC().Add(-1*monitor.alertWindow))

	requestCount := len(monitor.events) - index
	requestPerSecond := math.Ceil(float64(requestCount) / float64(5))

	monitor.cleanUp(index)

	if int(requestPerSecond) > monitor.threshold && !monitor.onAlert {
		monitor.onAlert = true
		monitor.out.WriteString(fmt.Sprintf("Monitor on alert %d requests/sec at time %s\n", int(requestPerSecond), time.Now().UTC().String()))
		// fmt.Printf( "Monitor on alert %.f requests/sec at time %s\n", requestPerSecond, time.Now().UTC().String())
	} else if int(requestPerSecond) < monitor.threshold && monitor.onAlert {
		monitor.onAlert = false
		monitor.out.WriteString(fmt.Sprintf("Monitor out of alert %d requests/sec at %s\n", int(requestPerSecond), time.Now().UTC().String()))
		// fmt.Printf("Monitor out of alert %.f requests/sec at %s\n", requestPerSecond, time.Now().UTC().String())
	}
}

func (monitor *alertMonitor) cleanUp(index int) {
	monitor.events = monitor.events[index:]
}

//// It will get the index of the first element that is before 2 minutes
func getIndexBeforeGivenTimeStamp(events []domain.Event, twoMinutes time.Time) int {
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
	return index
}
