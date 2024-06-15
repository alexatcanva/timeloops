package timeloops

import (
	"testing"
	"time"
)

func TestCountDuration(t *testing.T) {
	count := 5
	results := 0
	ForDuration(count, 5*time.Second, func(Break) {
		results++
	})
	if results != 5 {
		t.Fatal("did not reach expected count")
	}
}

func TestCountTimer(t *testing.T) {
	count := 5
	results := 0
	ForTimer(count, time.NewTimer(5*time.Second), func(Break) {
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

func TestBreakFuncFactoryNoN(t *testing.T) {
	testCases := []struct {
		desc   string
		n      int
		expect bool
	}{
		{
			desc:   "0 n",
			n:      0,
			expect: false,
		},
		{
			desc:   "-1 n",
			n:      -1,
			expect: false,
		},
		{
			desc:   "-1000 n",
			n:      -1000,
			expect: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := breakFuncFactory(tC.n)(tC.n)
			if res != tC.expect {
				t.Fatalf("expected: %t, got %t", tC.expect, res)
			}
		})
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
	executeForNIterationsOrTimeout(3, x, tc, fn)

	if ticks != 2 {
		t.Fatalf("expected: %d, got: %d", 2, ticks)
	}
}
