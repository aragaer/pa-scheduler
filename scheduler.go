package main

import (
	"container/list"
)

type Event struct {
	delay  int64
	repeat int64
	name   string
	what   []byte
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
		if queued.delay > event.delay {
			queued.delay -= event.delay
			scheduler.events.InsertBefore(event, e)
			return
		}
		event.delay -= queued.delay
	}
	scheduler.events.PushBack(event)
}

func (scheduler *Scheduler) Tick(seconds int) {
	first := scheduler.events.Front()
	if first != nil {
		first.Value.(*Event).delay -= int64(seconds)
	}
}

func (scheduler *Scheduler) GetTriggeredEvent() (triggered *Event) {
	first := scheduler.events.Front()
	if first == nil || first.Value.(*Event).delay > 0 {
		return
	}
	triggered = first.Value.(*Event)
	scheduler.events.Remove(first)
	if triggered.repeat != 0 {
		triggered.delay = triggered.repeat
		scheduler.Queue(triggered)
	}
	return
}
