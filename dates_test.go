package gengo

import (
	"testing"
	"time"
)

// DATE TESTS ----------------------------------------------------

func TestDate(t *testing.T) {
	minTime := time.Unix(dtMin, 0).UTC()
	maxTime := time.Unix(dtMax, 0).UTC()
	for i := 0; i < loop; i++ {
		d := Date()
		if d.Before(minTime) || d.After(maxTime) {
			t.Fatalf("Date() = %v: out of range [%v, %v]", d, minTime, maxTime)
		}
	}
}
func BenchmarkDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Date()
	}
}

// UNIX DATE TESTS -----------------------------------------------

func TestUnixDate(t *testing.T) {
	minTime := time.Unix(dtUnixMin, 0).UTC()
	maxTime := time.Unix(dtUnixMax, 0).UTC()
	for i := 0; i < loop; i++ {
		d := UnixDate()
		if d.Before(minTime) || d.After(maxTime) {
			t.Fatalf("UnixDate() = %v: out of range [%v, %v]", d, minTime, maxTime)
		}
	}
}
func BenchmarkUnixDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UnixDate()
	}
}

// DATE BETWEEN TESTS --------------------------------------------

func TestDateBetween(t *testing.T) {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC)
	for i := 0; i < loop; i++ {
		d := DateBetween(start, end)
		if d.Before(start) || d.After(end) {
			t.Fatalf("DateBetween(%v, %v) = %v: out of range", start, end, d)
		}
	}
	// swapped args must produce same valid range
	for i := 0; i < loop; i++ {
		d := DateBetween(end, start)
		if d.Before(start) || d.After(end) {
			t.Fatalf("DateBetween(end, start) = %v: out of range", d)
		}
	}
	// equal dates must return that exact time
	exact := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	if got := DateBetween(exact, exact); !got.Equal(exact) {
		t.Fatalf("DateBetween(t, t) = %v, want %v", got, exact)
	}
	// pre-epoch range: both boundaries are negative Unix timestamps
	preStart := time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC)
	preEnd := time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC)
	for i := 0; i < loop; i++ {
		d := DateBetween(preStart, preEnd)
		if d.Before(preStart) || d.After(preEnd) {
			t.Fatalf("DateBetween (pre-epoch) = %v: out of range [%v, %v]", d, preStart, preEnd)
		}
	}
}
func BenchmarkDateBetween(b *testing.B) {
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DateBetween(start, end)
	}
}
