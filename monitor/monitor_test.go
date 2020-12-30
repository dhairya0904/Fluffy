package monitor

import (
	"context"
	"fluffy/domain"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

/// It will start monitor for 15 seconds.
/// Report will be generated for the events sent
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

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Errorf("Error opening report file")
	}

	str := string(b)

	if len(str) == 0 {
		t.Errorf("Error, report file empty")
	}

	fmt.Println(str)
}
