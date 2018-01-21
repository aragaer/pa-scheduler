// -*- tab-width:4  -*-
package scheduler

import "encoding/json"

type scheduler struct {
	Events   <-chan string
	Commands chan<- []byte
	Ticks    chan<- int64
}

type cmd struct {
	Command string           `json:"command"`
	Delay   *int64           `json:"delay"`
	Repeat  *int64           `json:"repeat"`
	What    *json.RawMessage `json:"what"`
	Name    string           `json:"name"`
}

func New() *scheduler {
	cmdCh := make(chan []byte, 10)
	evtCh := make(chan string, 10)
	tickCh := make(chan int64, 10)
	evtQ := NewEventQueue()
	go start(cmdCh, tickCh, evtCh, evtQ)
	return &scheduler{evtCh, cmdCh, tickCh}
}

func (scheduler *scheduler) Close() {
	close(scheduler.Commands)
}

func (event *Event) fillFromCmd(cmd *cmd) {
	if cmd.What != nil {
		event.What = *cmd.What
	}
	if cmd.Delay != nil {
		event.Delay = *cmd.Delay
	}
	if cmd.Repeat != nil {
		event.Repeat = *cmd.Repeat
	}
}

func start(cmdCh <-chan []byte, tickCh <-chan int64, evtCh chan<- string, evtQ *eventQueue) {
Loop:
	for {
		select {
		case cmdBytes, ok := <-cmdCh:
			if ok == false {
				break Loop
			}
			var cmd cmd
			if err := json.Unmarshal(cmdBytes, &cmd); err == nil {
				switch cmd.Command {
				case "add":
					event := &Event{Name: cmd.Name}
					event.fillFromCmd(&cmd)
					evtQ.Queue(event)
				case "modify":
					event := evtQ.Remove(cmd.Name)
					event.fillFromCmd(&cmd)
					evtQ.Queue(event)
				case "cancel":
					evtQ.Remove(cmd.Name)
				}
			}
		case ticks, ok := <-tickCh:
			if ok == false {
				break Loop
			}
			evtQ.Tick(ticks)
		}
		for len(evtCh) < cap(evtCh) {
			event := evtQ.GetTriggeredEvent()
			if event == nil {
				break
			}
			what, _ := event.What.MarshalJSON()
			evtCh <- string(what)
		}
	}
	close(evtCh)
}
