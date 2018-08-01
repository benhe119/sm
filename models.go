package sm

import (
	"sync"
	"time"
)

// Event ...
type Event struct {
	sync.RWMutex

	ID            ID        `json:"id"`
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	RootCause     string    `json:"root_cause"`
	AffectedAreas []string  `json:"affected_areas"`
	Tags          []string  `json:"tags"`
	Level         int       `json:"level"`
	State         State     `json:"state"`
	CreatedAt     time.Time `json:"created"`
	MitigatedAt   time.Time `json:"mitigated"`
	FixedAt       time.Time `json:"fixed"`
	ClosedAt      time.Time `json:"closed"`
}

func NewEvent(title string, level int) (event *Event, err error) {
	event = &Event{
		ID:        db.NextId(),
		Title:     title,
		Level:     level,
		CreatedAt: time.Now(),
	}
	err = db.Save(event)
	if err == nil {
		metrics.Counter("event", "count").Inc()
	}
	return
}

func (e *Event) Id() ID {
	e.RLock()
	defer e.RUnlock()
	return e.ID
}

func (e *Event) Close() error {
	e.Lock()
	defer e.Unlock()
	e.State = STATE_CLOSED
	e.ClosedAt = time.Now()
	return db.Save(e)
}
