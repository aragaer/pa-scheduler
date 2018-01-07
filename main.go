package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func stdinLines(ch chan<- string) {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		ch <- scan.Text()
	}
	os.Exit(0)
}

func main() {
	ticker := time.Tick(1 * time.Second)
	last := time.Now().Unix()
	scheduler := NewScheduler()
	stdin := make(chan string)
	go stdinLines(stdin)
	for {
		select {
		case t := <-ticker:
			now := t.Unix()
			scheduler.Tick(now - last)
			for event := range scheduler.TriggeredEvents() {
				what, _ := event.What.MarshalJSON()
				fmt.Println(string(what))
			}
			last = now
		case line := <-stdin:
			event, err := Parse([]byte(line))
			if err == nil {
				scheduler.Queue(event)
			} else {
				fmt.Printf("Error parsing \"%s\": %v\n", line, err)
			}
		}
	}
}
