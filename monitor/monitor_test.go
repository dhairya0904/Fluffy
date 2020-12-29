package monitor

import (
	"fluffy/domain"
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	monitor := NewMonitor()

	events := make(chan domain.Event)
	done := make(chan bool)
	go monitor.StartMonitor(events, done)

	for _, event := range SyntheticData {
		events <- event
	}

	time.Sleep(15 * time.Second)
	done <- true
}
