package scheduler

import (
	"container/list"
	"encoding/json"
)

type Event struct {
	Delay  int64
	Repeat int64
	Name   string
	What   json.RawMessage
}

type eventQueue list.List

func NewEventQueue() (result *eventQueue) {
	return (*eventQueue)(list.New())
}

func (queue *eventQueue) asList() *list.List {
	return (*list.List)(queue)
}

func (queue *eventQueue) Front() (result *Event) {
	if e := queue.asList().Front(); e != nil {
		result = e.Value.(*Event)
	}
	return
}

func (queue *eventQueue) Queue(event *Event) {
	l := queue.asList()
	for e := l.Front(); e != nil; e = e.Next() {
		queued := e.Value.(*Event)
		if queued.Delay > event.Delay {
			queued.Delay -= event.Delay
			l.InsertBefore(event, e)
			return
		}
		event.Delay -= queued.Delay
	}
	l.PushBack(event)
}

func (queue *eventQueue) Tick(seconds int64) {
	if first := queue.Front(); first != nil {
		first.Delay -= seconds
	}
}

func (queue *eventQueue) GetTriggeredEvent() (result *Event) {
	first := queue.Front()
	if first == nil || first.Delay > 0 {
		return
	}
	result = first
	queue.Remove(result.Name)
	if result.Repeat != 0 {
		result.Delay %= result.Repeat
		result.Delay += result.Repeat
		queue.Queue(result)
	}
	return
}

func (queue *eventQueue) Add(event *Event) {
	for e := (*list.List)(queue).Front(); e != nil; e = e.Next() {
		if e.Value.(*Event).Name == event.Name {
			return
		}
	}
	queue.Queue(event)
}

func (queue *eventQueue) Remove(name string) (removed *Event) {
	for e := (*list.List)(queue).Front(); e != nil; e = e.Next() {
		queued := e.Value.(*Event)
		if queued.Name != name {
			continue
		}
		if e.Next() != nil {
			e.Next().Value.(*Event).Delay += queued.Delay
		}
		(*list.List)(queue).Remove(e)
		removed = queued
		break
	}
	return
}
