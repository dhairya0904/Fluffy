package alert

import (
	"fluffy/domain"
	"testing"
	"time"
)

func TestAlert(t *testing.T) {

	alertMonitor := NewAlertMonitor()
	events := make(chan domain.Event)
	done := make(chan bool)

	go alertMonitor.StartAlertMonitor(events, done)
	generateTraffic(events)
	if !alertMonitor.onAlert {
		t.Errorf("Monitor should be in alert state")
	}
	time.Sleep(5 * time.Second)

	if alertMonitor.onAlert {
		t.Errorf("Monitor should not be in alert state")
	}

}

func generateTraffic(events chan domain.Event) {

	timeout := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		timeout <- true
	}()

loop:
	for {
		select {
		default:
			events <- domain.Event{Time: time.Now().UTC()}
			time.Sleep(80 * time.Millisecond)
		case <-timeout:
			break loop
		}
	}
}
