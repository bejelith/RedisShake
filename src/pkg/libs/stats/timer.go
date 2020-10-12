package stats

import "time"

func NewTimer() *TimerContext {
	return &TimerContext{
		start: time.Now(),
	}
}

type TimerContext struct {
	start   time.Time
	elapsed time.Duration
}

func (t *TimerContext) Stop() time.Duration {
	if t.elapsed == 0 {
		t.elapsed = time.Now().Sub(t.start)
	}
	return t.elapsed
}

func (t *TimerContext) Duration() time.Duration {
	return t.Stop()
}
