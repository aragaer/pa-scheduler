# Scheduler

An application used for scheduling "events". Each event is a string that is printed to standard output.

Controlled by sending JSON commands to standard input

## Command format

Each command is a single-line JSON (no newlines must be within messages) ending with a newline. Each command must have the mandatory `command` field which determines which command to perform:

- `add`: add a new event

  If event with the same name already exists, command is silently discarded
  - `name`: unique identifier for a new event
  - `delay` _(optional)_: time in seconds until first occurrence of the event, default 0
  - `repeat` _(optional)_: time in seconds for periodic repeats or 0 for one-time event, default 0
  - `what`: string to be sent when event happens
- `modify`: change an existing event

  If event is not found, command is silently discarded
  - `name`: identifier of an event to be modified
  - `delay` _(optional)_: time in seconds until first occurrence of the event after the modification
  - `repeat` _(optional)_: time in seconds for periodic repeats or 0 for one-time event
  - `what` _(optional)_: if present, change the string to be sent when event happens
- `cancel`: remove a scheduled event

  If event is not found, command is silently discarded
  - `name`: identifier of an event to be cancelled
