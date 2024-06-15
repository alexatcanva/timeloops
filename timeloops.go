package timeloops

import (
	"time"
)

// Break is a function that will break a the timeloop.
type Break func()

// TimeLoopFunc is a function that will be executed in a timeloop.
type TimeloopFunc func(Break)

// ForDuration will perform some action in a loop for a given time.Duration, once
// the timer expires, or once n loops have been completed this function will
// return
func ForDuration(n int, duration time.Duration, fn TimeloopFunc) {
	bf := breakFuncFactory(n)
	timer := time.NewTimer(duration)
	executeForNIterationsOrTimeout(n, bf, timer.C, fn)
}

// ForTimer will perform some action in a loop for a given *time.Timer, once
// the timer expires, or once n loops have been completed this function will
// return
func ForTimer(n int, timer *time.Timer, fn TimeloopFunc) {
	executeForNIterationsOrTimeout(n, breakFuncFactory(n), timer.C, fn)
}

// breakFuncFactory will generate a function to break the loop
// dependening on whether or not a count (n) is provided.
// if n is <= 0, this is ignored.
func breakFuncFactory(n int) func(n int) bool {
	bf := func(n int) bool { return false }
	if n > 0 {
		bf = func(n int) bool {
			return n == 0
		}
	}
	return bf
}

// executeForNIterationsOrTimeout handles breaking the loop based off a read
// only time.Time chan.
var executeForNIterationsOrTimeout = func(
	n int,
	bfn func(n int) bool,
	stopChan <-chan time.Time,
	fn TimeloopFunc,
) {
	stop := false
	brk := func() {
		stop = true
	}
Loop:
	for !stop {
		if bfn(n) {
			break
		}
		select {
		case <-stopChan:
			break Loop
		default:
			fn(brk)
			n--
		}
	}
}
