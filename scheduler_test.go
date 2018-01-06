package main

import (
	"strings"
	"testing"
)

type event struct {
	delay, repeat int64
	name          string
}

func (e event) mk() *Event {
	return &Event{e.delay, e.repeat, e.name, []byte(`{"event": "tick"}`)}
}

type Events []event
type Expected []string

type tc struct {
	events   Events
	expected Expected
}

var testCases = map[string]tc{
	"one immediate": {
		Events{{0, 0, "tick"}},
		Expected{"tick", "", ""}},
	"one delayed": {
		Events{{1, 0, "tick"}},
		Expected{"", "tick", ""}},
	"one repeating": {
		Events{{0, 1, "tick"}},
		Expected{"tick", "tick", "tick"}},
	"two immediate": {
		Events{{0, 0, "tick"}, {0, 0, "tock"}},
		Expected{"tick tock", "", ""}},
	"two alternating": {
		Events{{0, 2, "tick"}, {1, 2, "tock"}},
		Expected{"tick", "tock", "tick", "tock"}},
	"two different freq": {
		Events{{0, 2, "tick"}, {0, 3, "tock"}},
		Expected{"tick tock", "", "tick", "tock"}},
	"insert before": {
		Events{{1, 0, "tock"}, {0, 0, "tick"}},
		Expected{"tick", "tock", ""}},
}

func TestScheduler(t *testing.T) {
	for name, tc := range testCases {
		scheduler := NewScheduler()
		for _, e := range tc.events {
			scheduler.Queue(e.mk())
		}

		for tick, eventsForTick := range tc.expected {
			for _, expected := range strings.Fields(eventsForTick) {
				event := scheduler.GetTriggeredEvent()
				if event == nil {
					t.Errorf("Test case \"%s\" failed:", name)
					t.Fatalf("\"%s\" event is not triggered on tick %d", expected, tick)
				} else if event.Name != expected {
					t.Errorf("Test case \"%s\" failed:", name)
					t.Fatalf("Expected \"%s\", got \"%s\"", expected, event.Name)
				}
			}
			event := scheduler.GetTriggeredEvent()
			if event != nil {
				t.Errorf("Test case \"%s\" failed:", name)
				t.Fatalf("\"%s\" event has ticked on tick %d", event.Name, tick)
			}

			scheduler.Tick(1)
		}
	}
}
