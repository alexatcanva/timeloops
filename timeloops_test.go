package timeloop

import (
	"testing"
	"time"
)

func TestCount(t *testing.T) {
	count := 5
	results := 0
	ForDuration(count, 5*time.Second, func() {
		results += 1
	})
	if results != 5 {
		t.Fatal("did not reach expected count")
	}
}

func TestBreakFuncFactory(t *testing.T) {
	n := 3
	bf := breakFuncFactory(n)
	count := 0
	for !bf(n) {
		count++
		n--
	}
	if count != 3 {
		t.Fatalf("expected %d, got %d", 3, count)
	}
}

func TestForDuringChan(t *testing.T) {
	tc := make(chan time.Time, 1)
	x := func(n int) bool { return false }

	ticks := 0
	fn := func() {
		ticks++
		if ticks == 2 {
			tc <- time.Now()
		}
	}
	forDurationChan(3, x, tc, fn)

	if ticks != 2 {
		t.Fatalf("expected: %d, got: %d", 2, ticks)
	}
}
