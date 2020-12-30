package alert

import (
	"context"
	"fluffy/domain"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

/// Generator will generate traffic of 10 req/s
/// monitor will be in alert state after 5 seconds
func TestAlert(t *testing.T) {

	alertMonitor := NewAlertMonitor(5, 10)
	events := make(chan domain.Event)
	ctx, cancel := context.WithCancel(context.Background())
	f, _ := ioutil.TempFile("", "access")
	defer os.Remove(f.Name())
	var wg sync.WaitGroup
	wg.Add(1)

	go alertMonitor.StartAlertMonitor(ctx, &wg, events, f.Name())
	generateTraffic(events, 10)
	if !alertMonitor.onAlert {
		t.Errorf("Monitor should be in alert state")
	}
	time.Sleep(5 * time.Second)

	if alertMonitor.onAlert {
		t.Errorf("Monitor should not be in alert state")
	}
	cancel()
}

/// This will generate traffic for only 4 seconds
/// requests/sec is 10 but it is out of window
func TestAlertForTrafficNotInWindow(t *testing.T) {

	alertMonitor := NewAlertMonitor(5, 10)
	events := make(chan domain.Event)
	ctx, cancel := context.WithCancel(context.Background())
	f, _ := ioutil.TempFile("", "access")
	defer os.Remove(f.Name())
	var wg sync.WaitGroup
	wg.Add(1)

	go alertMonitor.StartAlertMonitor(ctx, &wg, events, f.Name())
	generateTraffic(events, 4)
	if alertMonitor.onAlert {
		t.Errorf("Monitor should be in alert state")
	}
	cancel()
}

/// This will generate traffic for 10 seconds
/// One request in 80 ms, which will generate a traffic of above 10 requests per second
func generateTraffic(events chan domain.Event, duration int) {

	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
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
