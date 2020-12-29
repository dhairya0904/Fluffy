package reader

import (
	"bufio"
	"fluffy/domain"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const SyntheticData = "127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"

func TestReader(t *testing.T) {

	reader := NewReader()
	subs1 := make(chan domain.Event)
	subs2 := make(chan domain.Event)

	reader.Subscribe(subs1)
	reader.Subscribe(subs2)

	f, _ := ioutil.TempFile("", "access")
	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)

	doneChan := make(chan bool)
	go reader.StartPublishing(f.Name(), doneChan)
	w.WriteString(SyntheticData + "\n")

	go func() {
		time.Sleep(2 * time.Second)
		doneChan <- true
	}()

	for event := range subs1 {
		if !assertEvent(event) {
			t.Errorf("Error in parsing")
		}
	}
	for event := range subs2 {
		if !assertEvent(event) {
			t.Errorf("Error in parsing")
		}
	}
}

func assertEvent(event domain.Event) bool {
	if event.Host != "127.0.0.1" && event.User != "james" && event.Method != "GET" &&
		event.URL != "/report" && event.Bytes != 123 && event.Status != 200 &&
		event.Protocol != "HTTP/1.0" {
		return false
	}
	return true
}
