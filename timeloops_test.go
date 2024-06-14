package timeloop

import (
	"math"
	"testing"
	"time"
)

func TestCount(t *testing.T) {
	count := 5
	results := 0
	ForDuration(count, 5*time.Second, func() {
		print("hello")
		results += 1
	})

	if results != 5 {
		t.Fatal("did not reach expected count")
	}
}

func TestTimeNaive(t *testing.T) {
	before := time.Now()
	ForDuration(0, 5*time.Second, func() {
		print("hello")
	})
	after := time.Now()
	got := math.Floor(after.Sub(before).Seconds())
	if got != 5 {
		t.Fatalf("expected %f, got %f", 5.0, got)
	}
}
