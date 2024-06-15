package timeloops

import (
	"errors"
	"time"
)

// Return TimeloopBreak to break the timeloop.
//
//lint:ignore ST1012 Overload the error type for easy API.
var TimeloopBreak error = errors.New("break")

// ForDuration will perform some action in a loop for a given time.Duration,
// once the timer expires, or once n loops have been completed this function
// will return
func ForDuration(n int, duration time.Duration, fn func() error) error {
	bf := breakFuncFactory(n)
	timer := time.NewTimer(duration)
	return executeForNIterationsOrTimeout(n, bf, timer.C, fn)
}

// ForTimer will perform some action in a loop for a given *time.Timer, once
// the timer expires, or once n loops have been completed this function will
// return
func ForTimer(n int, timer *time.Timer, fn func() error) error {
	return executeForNIterationsOrTimeout(n, breakFuncFactory(n), timer.C, fn)
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
	fn func() error,
) error {
Loop:
	for {
		if bfn(n) {
			break
		}
		select {
		case <-stopChan:
			break Loop
		default:
			if err := fn(); err != nil {
				if errors.Is(err, TimeloopBreak) {
					break Loop
				}
				return err
			}
			n--
		}
	}
	return nil
}
