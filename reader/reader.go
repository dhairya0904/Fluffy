package reader

import (
	"context"
	"fluffy/domain"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hpcloud/tail"
	"xojoc.pw/logparse"
)

//// Reader: It will read from logs and publish event to all subscribed channels
type Reader struct {
	subs []chan domain.Event
}

func NewReader() *Reader {
	reader := &Reader{}
	reader.subs = make([]chan domain.Event, 0)
	return reader
}

func (reader *Reader) Subscribe(c chan domain.Event) {
	reader.subs = append(reader.subs, c)
}

//// Starts reading file and publishing events to all
//// the subscribed channels
func (reader *Reader) StartPublishing(ctx context.Context, wg *sync.WaitGroup, filename string) {
	t, err := tail.TailFile(filename, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}})
	defer wg.Done()

	if err != nil {
		log.Print("Error while opening file")
		os.Exit(1)
	}

	log.Println("starting monitor")

	go func() {
		<-ctx.Done()
		t.Stop()
		for _, subs := range reader.subs {
			close(subs)
		}
	}()

	for line := range t.Lines {
		for _, ch := range reader.subs {
			msg, err := parseLogs(line.Text)
			if err != nil {
				log.Println("Error parsing logs for", line.Text, err)
				continue
			}
			go func(ch chan domain.Event) {
				ch <- msg
			}(ch)
		}
	}
}

func parseLogs(msg string) (domain.Event, error) {

	l, err := logparse.Common(msg)

	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		Host:     l.Host.String(),
		User:     l.User,
		Time:     l.Time,
		Method:   l.Request.Method,
		URL:      l.Request.URL.Path,
		Protocol: l.Request.Proto,
		Status:   l.Status,
		Bytes:    l.Bytes,
	}, err
}
