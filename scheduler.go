package main

import (
	"container/list"
	"encoding/json"
	"errors"
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

func (scheduler *Scheduler) Tick(seconds int64) {
	if scheduler.events.Len() > 0 {
		scheduler.events.Front().Value.(*Event).Delay -= seconds
	}
}

func (scheduler *Scheduler) putTriggeredEventsToChannel(ch chan<- *Event) {
	for {
		first := scheduler.events.Front()
		if first == nil || first.Value.(*Event).Delay > 0 {
			break
		}
		triggered := first.Value.(*Event)
		ch <- triggered
		scheduler.events.Remove(first)
		if triggered.Repeat != 0 {
			triggered.Delay = triggered.Repeat
			scheduler.Queue(triggered)
		}
	}
	close(ch)
}

func (scheduler *Scheduler) TriggeredEvents() <-chan *Event {
	ch := make(chan *Event)
	go scheduler.putTriggeredEventsToChannel(ch)
	return ch
}

func Parse(message []byte) (result *Event, err error) {
	err = json.Unmarshal(message, &result)
	if err == nil && result.Name == "" {
		result = nil
		err = errors.New("\"name\" field is missing")
	}
	return
}

func (scheduler *Scheduler) Add(event *Event) {
	for e := scheduler.events.Front(); e != nil; e = e.Next() {
		if e.Value.(*Event).Name == event.Name {
			return
		}
	}
	scheduler.Queue(event)
}

func (scheduler *Scheduler) Modify(event *Event) {
	old := scheduler.Remove(event)
	if old != nil {
		old.Delay = event.Delay
		old.Repeat = event.Repeat
		if event.What != nil {
			old.What = event.What
		}
		scheduler.Queue(old)
	}
}

func (scheduler *Scheduler) Remove(event *Event) (removed *Event) {
	for e := scheduler.events.Front(); e != nil; e = e.Next() {
		queued := e.Value.(*Event)
		if queued.Name == event.Name {
			if e.Next() != nil {
				e.Next().Value.(*Event).Delay += queued.Delay
			}
			scheduler.events.Remove(e)
			removed = queued
			break
		}
	}
	return
}
