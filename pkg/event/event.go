package event

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
)

const (
	EventReceivedName = "event/event.received"
	FnFailedName      = "inngest/function.failed"
)

var (
	startTimestamp = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	endTimestamp   = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
)

type TrackedEvent interface {
	GetInternalID() ulid.ULID
	GetEvent() Event
}

func NewEvent(data string) (*Event, error) {
	evt := &Event{}
	if err := json.Unmarshal([]byte(data), evt); err != nil {
		return nil, err
	}

	return evt, nil
}

// Event represents an event sent to Inngest.
type Event struct {
	Name string         `json:"name"`
	Data map[string]any `json:"data"`

	// User represents user-specific information for the event.
	User map[string]any `json:"user,omitempty"`

	// ID represents the unique ID for this particular event.  If supplied, we should attempt
	// to only ingest this event once.
	ID string `json:"id,omitempty"`

	// Timestamp is the time the event occurred, at millisecond precision.
	// If this is not provided, we will insert the current time upon receipt of the event
	Timestamp int64  `json:"ts,omitempty"`
	Version   string `json:"v,omitempty"`
}

func (evt Event) Time() time.Time {
	return time.UnixMilli(evt.Timestamp)
}

func (evt Event) Map() map[string]any {
	if evt.Data == nil {
		evt.Data = make(map[string]any)
	}
	if evt.User == nil {
		evt.User = make(map[string]any)
	}

	data := map[string]any{
		"name": evt.Name,
		"data": evt.Data,
		"user": evt.User,
		"id":   evt.ID,
		// We cast to float64 because marshalling and unmarshalling from
		// JSON automatically uses float64 as its type;  JS has no notion
		// of ints.
		"ts": float64(evt.Timestamp),
	}

	if evt.Version != "" {
		data["v"] = evt.Version
	}

	return data
}

func (e Event) Validate(ctx context.Context) error {
	if e.Name == "" {
		return errors.New("event name is empty")
	}

	if e.Timestamp != 0 {
		// Convert milliseconds to nanosecond precision
		t := time.Unix(0, e.Timestamp*1_000_000)
		if t.Before(startTimestamp) {
			return errors.New("timestamp is before Jan 1, 1980")
		}
		if t.After(endTimestamp) {
			return errors.New("timestamp is after Jan 1, 2100")
		}
	}

	return nil
}

func NewOSSTrackedEvent(e Event) TrackedEvent {
	id, err := ulid.Parse(e.ID)
	if err != nil {
		id = ulid.MustNew(ulid.Now(), rand.Reader)
	}
	if e.ID == "" {
		e.ID = id.String()
	}
	return ossTrackedEvent{
		id:    id,
		event: e,
	}
}

type ossTrackedEvent struct {
	id    ulid.ULID
	event Event
}

func (o ossTrackedEvent) GetEvent() Event {
	return o.event
}

func (o ossTrackedEvent) GetInternalID() ulid.ULID {
	return o.id
}
