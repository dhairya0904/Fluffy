package main

import (
	"fmt"
	"os"
	"time"
)

const data = "127.0.0.1 - james [%s] \"GET /report HTTP/1.0\" 200 123\n"

func main() {
	f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	timeout := make(chan bool)
	go func() {
		time.Sleep(15 * time.Second)
		timeout <- true
	}()

loop:
	for {
		select {
		default:
			f.WriteString(fmt.Sprintf(data, time.Now().UTC().Format("02/Jan/2006:15:04:05 -0700")))
			time.Sleep(10 * time.Millisecond)
		case <-timeout:
			break loop
		}
	}
}
