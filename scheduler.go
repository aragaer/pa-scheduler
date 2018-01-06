package main

import (
	"container/list"
	"encoding/json"
)

type Event struct {
	Delay  int64           `json:"delay"`
	Repeat int64           `json:"repeat"`
	Name   string          `json:"name"`
	What   json.RawMessage `json:"what"`
}

type Scheduler struct {
	events *list.List
}

func NewScheduler() (result *Scheduler) {
	return &Scheduler{list.New()}
}

func (scheduler *Scheduler) Queue(event *Event) {
	for e := scheduler.events.Front(); e != nil; e = e.Next() {
		queued := e.Value.(*Event)
		if queued.Delay > event.Delay {
			queued.Delay -= event.Delay
			scheduler.events.InsertBefore(event, e)
			return
		}
		event.Delay -= queued.Delay
	}
	scheduler.events.PushBack(event)
}

func (scheduler *Scheduler) Tick(seconds int) {
	first := scheduler.events.Front()
	if first != nil {
		first.Value.(*Event).Delay -= int64(seconds)
	}
}

func (scheduler *Scheduler) GetTriggeredEvent() (triggered *Event) {
	first := scheduler.events.Front()
	if first == nil || first.Value.(*Event).Delay > 0 {
		return
	}
	triggered = first.Value.(*Event)
	scheduler.events.Remove(first)
	if triggered.Repeat != 0 {
		triggered.Delay = triggered.Repeat
		scheduler.Queue(triggered)
	}
	return
}

func Parse(message []byte) (result *Event, err error) {
	event := Event{}
	err = json.Unmarshal(message, &event)
	return &event, err
}
