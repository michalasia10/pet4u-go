package clock

import "time"

// Port defines a time source to allow deterministic testing.
type Port interface {
    Now() time.Time
}


