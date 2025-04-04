package gengo

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	for i := 0; i < 1000; i++ {
		Date()
	}
}

func TestUnixDate(t *testing.T) {
	for i := 0; i < 1000; i++ {
		UnixDate()
	}
}

func TestDateBetween(t *testing.T) {

	time.Now().Unix()
	dMin := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	dMax := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	for i := 0; i < 1000; i++ {
		DateBetween(dMin, dMax)
	}
}
