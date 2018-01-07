package main

import (
	"bufio"
	"encoding/json"
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
			var command map[string]interface{}
			bytes := []byte(line)
			if err := json.Unmarshal(bytes, &command); err == nil {
				action := command["command"].(string)
				event, err := Parse(bytes)
				if err == nil {
					switch action {
					case "add":
						scheduler.Add(event)
					case "modify":
						scheduler.Modify(event)
					case "cancel":
						scheduler.Remove(event)
					}
				} else {
					fmt.Printf("Error parsing \"%s\": %v\n", line, err)
				}
			}
		}
	}
}
