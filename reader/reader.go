package reader

import (
	"fluffy/domain"
	"io"
	"log"

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
func (reader *Reader) StartPublishing(filename string, doneChan chan bool) {
	t, _ := tail.TailFile(filename, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}})

	go func() {
		<-doneChan
		t.Stop()
		for _, subs := range reader.subs {
			close(subs)
		}
	}()

	for line := range t.Lines {
		for _, ch := range reader.subs {
			msg, err := parseLogs(line.Text)
			if err != nil {
				log.Print("Error parsing logs for", line.Text)
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
		log.Print("Error parsing logs")
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
