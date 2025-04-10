package gengo

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	for i := 0; i < loop; i++ {
		Date()
	}
}
func BenchmarkDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Date()
	}
}

func TestUnixDate(t *testing.T) {
	for i := 0; i < loop; i++ {
		UnixDate()
	}
}
func BenchmarkUnixDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UnixDate()
	}
}

func TestDateBetween(t *testing.T) {

	dMin := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	dMax := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	for i := 0; i < loop; i++ {
		DateBetween(dMin, dMax)
	}
}
func BenchmarkDateBetween(b *testing.B) {
	dMin := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	dMax := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DateBetween(dMin, dMax)
	}
}
