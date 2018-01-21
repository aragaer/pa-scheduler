package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"github.com/aragaer/scheduler"
)

func stdinLines(ch chan<- string) {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		ch <- scan.Text()
	}
	close(ch)
}

func main() {
	ticker := time.Tick(time.Second / 2)
	last := time.Now().Unix()
	sched := scheduler.New()
	stdin := make(chan string)
	go stdinLines(stdin)
	for {
		select {
		case t := <-ticker:
			now := t.Unix()
			sched.Ticks <- now - last
			last = now
		case evt := <- sched.Events:
			fmt.Println(evt)
		case line, ok := <-stdin:
			if !ok {
				return
			}
			sched.Commands <- []byte(line)
		}
	}
}
