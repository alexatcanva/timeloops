package timeloop

import (
	"time"
)

// ForDuration is a function that runs a function for a specified number of
// times or until a specified time has passed.
func ForDuration(n int, duration time.Duration, fn func()) {
	var bf func(n int) bool = func(n int) bool { return false }
	if n > 0 {
		bf = func(n int) bool {
			return n == 0
		}
	}
	timer := time.NewTimer(duration)
Loop:
	for {
		if bf(n) {
			break
		}
		select {
		case <-timer.C:
			break Loop
		default:
			fn()
			n--
		}
	}
}
