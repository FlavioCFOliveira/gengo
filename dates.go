package gengo

import (
	"math/rand/v2"
	"time"
)

func Date() time.Time {
	randomUnix := rand.Int64N(dtMax-dtMin+1) + dtMin
	return time.Unix(randomUnix, 0).UTC()
}
func UnixDate() time.Time {
	randomUnix := rand.Int64N(dtUnixMax-dtUnixMin+1) + dtUnixMin
	return time.Unix(randomUnix, 0).UTC()
}

func DateBetween(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}

	startUnix := start.Unix()
	endUnix := end.Unix()

	if startUnix == endUnix {
		return start
	}

	// Gera um número aleatório entre startUnix e endUnix
	randomUnix := rand.Int64N(endUnix-startUnix+1) + startUnix

	// Converte para time.Time
	return time.Unix(randomUnix, 0).UTC()
}
