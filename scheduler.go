package main

type Event struct {
	delay  int64
	repeat int64
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

func (scheduler *scheduler) Tick(seconds int) {
	if scheduler.event != nil {
		scheduler.event.delay -= int64(seconds)
	}
}

func (scheduler *scheduler) GetTriggeredEvent() *Event {
	if scheduler.event == nil || scheduler.event.delay > 0 {
		return nil
	}
	result := scheduler.event
	scheduler.event = nil
	return result
}
