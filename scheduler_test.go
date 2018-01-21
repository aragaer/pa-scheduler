package scheduler

import "time"
import "testing"

func Verify(expected []string, scheduler *scheduler, t *testing.T) {
	for {
		select {
		case message := <-scheduler.Events:
			if len(expected) == 0 {
				t.Fatalf("Expected nothing, got event '%s'", message)
			}
			if message != expected[0] {
				t.Fatalf("Expected '%s', got event '%s'", expected[0], message)
			}
			expected = expected[1:]
		case <-time.After(time.Millisecond):
			if len(expected) > 0 {
				t.Fatalf("Expected '%s', got nothing", expected[0])
			}
			return
		}
	}
}

func TestSchedulerAdd(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()
	Verify([]string{}, scheduler, t)

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello"}`)
	Verify([]string{`"hello"`}, scheduler, t)

	scheduler.Ticks <- 1
	Verify([]string{}, scheduler, t)

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello2", "delay": 3}`)
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 2
	Verify([]string{`"hello2"`}, scheduler, t)

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello3", "repeat": 1}`)
	Verify([]string{`"hello3"`}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{`"hello3"`}, scheduler, t)
	scheduler.Ticks <- 2
	Verify([]string{`"hello3"`}, scheduler, t)
}

func TestSchedulerModify1(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello", "delay": 100}`)
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{}, scheduler, t)

	scheduler.Commands <- []byte(`{"command": "modify", "name": "hello", "delay": 1}`)
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{`"hello"`}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{}, scheduler, t)
}

func TestSchedulerModify2(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello", "repeat": 1}`)
	Verify([]string{`"hello"`}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{`"hello"`}, scheduler, t)

	scheduler.Commands <- []byte(`{"command": "modify", "name": "hello", "what": "hi"}`)
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{`"hi"`}, scheduler, t)
	scheduler.Ticks <- 1
	Verify([]string{`"hi"`}, scheduler, t)
}

func TestSchedulerModify3(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello", "repeat": 2}`)
	Verify([]string{`"hello"`}, scheduler, t)
	scheduler.Ticks <- 1

	scheduler.Commands <- []byte(`{"command": "modify", "name": "hello", "repeat": 3}`)
	// Wait until scheduler handles the command
	for len(scheduler.Commands) > 0 {
		time.Sleep(time.Millisecond)
	}
	scheduler.Ticks <- 1
	Verify([]string{`"hello"`}, scheduler, t)
	scheduler.Ticks <- 2
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 2
	Verify([]string{`"hello"`}, scheduler, t)
}

func TestSchedulerCancel(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello", "repeat": 2}`)
	Verify([]string{`"hello"`}, scheduler, t)
	scheduler.Ticks <- 1

	scheduler.Commands <- []byte(`{"command": "cancel", "name": "hello"}`)
	// Wait until scheduler handles the command
	for len(scheduler.Commands) > 0 {
		time.Sleep(time.Millisecond)
	}
	scheduler.Ticks <- 1
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 2
	Verify([]string{}, scheduler, t)
	scheduler.Ticks <- 2
	Verify([]string{}, scheduler, t)
}

func TestSchedulerTwoEvents(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello1", "what": "hello1"}`)
	scheduler.Commands <- []byte(`{"command": "add", "name": "hello2", "what": "hello2", "delay": 1}`)
	Verify([]string{`"hello1"`}, scheduler, t)

	scheduler.Ticks <- 1
	Verify([]string{`"hello2"`}, scheduler, t)
}

func TestSchedulerCancelTwo(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello1", "what": "hello1", "delay": 10}`)
	scheduler.Commands <- []byte(`{"command": "add", "name": "hello2", "what": "hello2", "delay": 1}`)
	scheduler.Commands <- []byte(`{"command": "cancel", "name": "hello1"}`)
	// Wait until scheduler handles the command
	for len(scheduler.Commands) > 0 {
		time.Sleep(time.Millisecond)
	}

	scheduler.Ticks <- 20
	Verify([]string{`"hello2"`}, scheduler, t)
}

func TestSchedulerLock(t *testing.T) {
	scheduler := New()
	defer scheduler.Close()

	scheduler.Commands <- []byte(`{"command": "add", "name": "hello", "what": "hello", "repeat": 1}`)
	for i := 0; i < 100; i++ {
		scheduler.Ticks <- 1
	}
}

func TestSchedulerCleanup(t *testing.T) {
	scheduler := New()

	scheduler.Close()

	if _, ok := <-scheduler.Events; ok {
		t.Fatalf("Scheduler is not finalized")
	}
}

func TestSchedulerCleanup2(t *testing.T) {
	scheduler := New()

	close(scheduler.Ticks)

	if _, ok := <-scheduler.Events; ok {
		t.Fatalf("Scheduler is not finalized")
	}
}
