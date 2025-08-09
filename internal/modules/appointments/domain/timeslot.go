package domain

import (
    "errors"
    "time"
)

var ErrInvalidTimeSlot = errors.New("invalid time slot: start must be before end")

// TimeSlot is a value object with an invariant Start < End.
type TimeSlot struct {
    Start time.Time
    End   time.Time
}

func NewTimeSlot(start, end time.Time) (TimeSlot, error) {
    if !start.Before(end) {
        return TimeSlot{}, ErrInvalidTimeSlot
    }
    return TimeSlot{Start: start, End: end}, nil
}

func (s TimeSlot) Duration() time.Duration { return s.End.Sub(s.Start) }

func (s TimeSlot) Overlaps(other TimeSlot) bool {
    return s.Start.Before(other.End) && other.Start.Before(s.End)
}


