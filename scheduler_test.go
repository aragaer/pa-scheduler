package main

import (
	"testing"
)

func TestScheduler_ImmediateOnce(t *testing.T) {
	scheduler := NewScheduler()

	if scheduler == nil {
		t.Fatalf("Scheduler is not created")
	}

	scheduler.Queue(&Event{
		delay:  0,
		repeat: 0,
		name:   "tick once to brain1",
		what:   []byte(`{"event": "tick", "to": "brain1"}`),
	})

	event := scheduler.GetTriggeredEvent()

	if event == nil {
		t.Fatalf("\"tick once\" event is not triggered")
	} else if event.name != "tick once to brain1" {
		t.Fatalf("Expected \"tick once to brain1\", got \"%s\"", event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}

func TestScheduler_DelayedOnce(t *testing.T) {
	scheduler := NewScheduler()

	if scheduler == nil {
		t.Fatalf("Scheduler is not created")
	}

	scheduler.Queue(&Event{
		delay:  1,
		repeat: 0,
		name:   "tick once to brain1",
		what:   []byte(`{"event": "tick", "to": "brain1"}`),
	})

	event := scheduler.GetTriggeredEvent()

	if event != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}

	scheduler.Tick(1)

	event = scheduler.GetTriggeredEvent()

	if event == nil {
		t.Fatalf("\"tick once\" event is not triggered")
	} else if event.name != "tick once to brain1" {
		t.Fatalf("Expected \"tick once to brain1\", got \"%s\"", event.name)
	}

	if scheduler.GetTriggeredEvent() != nil {
		t.Fatalf("\"%s\" event has ticked", event.name)
	}
}
