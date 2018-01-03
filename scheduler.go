package main

type Event struct {
	delay  uint64
	repeat uint64
	name   string
	what   []byte
}

type scheduler struct {
	event *Event
}

func NewScheduler() *scheduler {
	return &scheduler{}
}

func (scheduler *scheduler) Queue(event *Event) {
	scheduler.event = event
}

func (scheduler *scheduler) GetTriggeredEvent() *Event {
	result := scheduler.event
	scheduler.event = nil
	return result
}
