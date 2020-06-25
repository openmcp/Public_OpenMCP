package clock

import "time"

var MyClock clock = &realClock{}

type realClock struct{}

func (realClock) Now() time.Time                  { return time.Now() }
func (realClock) Since(d time.Time) time.Duration { return time.Since(d) }

type clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}
