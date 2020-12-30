package monitor

import (
	"context"
	"fluffy/domain"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	monitor := NewMonitor(10)

	events := make(chan domain.Event)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	f, _ := ioutil.TempFile("", "report")
	defer os.Remove(f.Name())
	wg.Add(1)

	go monitor.StartMonitor(ctx, &wg, events, f.Name())

	for _, event := range SyntheticData {
		events <- event
	}

	time.Sleep(15 * time.Second)
	cancel()
	wg.Wait()
}
