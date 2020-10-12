package stats

import (
	"testing"
	"time"
)

// Use tens of second to approximate results
var TenthOfSecond = time.Second * 10
var delays = []time.Duration{time.Second / 2, time.Second / 4}

func TestTimerMeasuresIsCorrect(t *testing.T) {
	for _, delay := range delays {
		timer := NewTimer()
		time.Sleep(delay)
		if timer.Stop()/TenthOfSecond != delay/TenthOfSecond {
			t.Fatalf("Elapsed: %d, wanted %d", timer.Duration()/TenthOfSecond, delay/TenthOfSecond)
		}
	}
}

func TestTimerStopEqualsDuration(t *testing.T) {
	for _, delay := range delays {
		timer := NewTimer()
		time.Sleep(delay)
		if timer.Stop() != timer.Duration() {
			t.Fatalf("Stop: %d, duration %d", timer.Duration(), timer.Duration())
		}
	}
}
