package main

type Event struct {
	delay  int64
	repeat int64
	name   string
	what   []byte
}

type Scheduler struct {
	event *Event
}

func (scheduler *Scheduler) Queue(event *Event) {
	scheduler.event = event
}

func (scheduler *Scheduler) Tick(seconds int) {
	if scheduler.event != nil {
		scheduler.event.delay -= int64(seconds)
	}
}

func (scheduler *Scheduler) GetTriggeredEvent() *Event {
	if scheduler.event == nil || scheduler.event.delay > 0 {
		return nil
	}
	result := scheduler.event
	if result.repeat != 0 {
		result.delay = result.repeat
	} else {
		scheduler.event = nil
	}
	return result
}
