package main

import (
	"testing"
)

func MakeEvent(delay int64, repeat int64, name string) *Event {
	return &Event{
		delay:  delay,
		repeat: repeat,
		name:   name,
		what:   []byte(`{"event": "tick", "to": "brain1"}`),
	}
}

func TestScheduler_ImmediateOnce(t *testing.T) {
	scheduler := NewScheduler()
	name := "tick once to brain1"

	scheduler.Queue(MakeEvent(0, 0, name))

	event := scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name)
	} else if event.name != name {
		t.Fatalf("Expected \"%s\", got \"%s\"", name, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}

func TestScheduler_DelayedOnce(t *testing.T) {
	scheduler := NewScheduler()
	name := "tick once to brain1"

	scheduler.Queue(MakeEvent(1, 0, name))

	event := scheduler.GetTriggeredEvent()
	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}

	scheduler.Tick(1)

	event = scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name)
	} else if event.name != name {
		t.Fatalf("Expected \"%s\", got \"%s\"", name, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}

func TestScheduler_Repeating(t *testing.T) {
	scheduler := NewScheduler()
	name := "tick to brain1"

	scheduler.Queue(MakeEvent(0, 1, name))

	event := scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name)
	} else if event.name != name {
		t.Fatalf("Expected \"%s\", got \"%s\"", name, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}

	scheduler.Tick(1)

	event = scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not repeated", name)
	} else if event.name != name {
		t.Fatalf("Expected \"%s\", got \"%s\"", name, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}

func TestScheduler_TwoImmediateOnce(t *testing.T) {
	scheduler := NewScheduler()
	name1 := "tick1"
	name2 := "tick2"

	scheduler.Queue(MakeEvent(0, 0, name1))
	scheduler.Queue(MakeEvent(0, 0, name2))

	event := scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name1)
	} else if event.name != name1 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name1, event.name)
	}

	event = scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name2)
	} else if event.name != name2 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name2, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}

func TestScheduler_TwoAlternateRepeating(t *testing.T) {
	scheduler := NewScheduler()
	name1 := "tick1"
	name2 := "tick2"

	scheduler.Queue(MakeEvent(0, 2, name1))
	scheduler.Queue(MakeEvent(1, 2, name2))

	event := scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name1)
	} else if event.name != name1 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name1, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}

	scheduler.Tick(1)

	event = scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name2)
	} else if event.name != name2 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name2, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}

	scheduler.Tick(1)

	event = scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name1)
	} else if event.name != name1 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name1, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}

func TestScheduler_ScheduleBefore(t *testing.T) {
	scheduler := NewScheduler()
	name1 := "tick1"
	name2 := "tick2"

	scheduler.Queue(MakeEvent(1, 0, name2))
	scheduler.Queue(MakeEvent(0, 0, name1))

	event := scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name1)
	} else if event.name != name1 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name1, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}

	scheduler.Tick(1)

	event = scheduler.GetTriggeredEvent()
	if event == nil {
		t.Fatalf("\"%s\" event is not triggered", name2)
	} else if event.name != name2 {
		t.Fatalf("Expected \"%s\", got \"%s\"", name2, event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}
