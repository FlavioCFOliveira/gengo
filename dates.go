package gengo

import (
	"math/rand/v2"
	"time"
)

// Date returns a random time.Time uniformly distributed across the full Gregorian calendar range (year 1–9999).
func Date() time.Time {
	return time.Unix(rand.Int64N(dtRange)+dtMin, 0).UTC()
}

// UnixDate returns a random time.Time uniformly distributed within the Unix epoch range (1970–2038).
func UnixDate() time.Time {
	return time.Unix(rand.Int64N(dtUnixRange)+dtUnixMin, 0).UTC()
}

// DateBetween returns a random time.Time in [start, end]. Arguments are swapped if start is after end.
func DateBetween(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}

	startUnix := start.Unix()
	endUnix := end.Unix()

	if startUnix == endUnix {
		return start
	}

	// Generate a random number between startUnix and endUnix
	randomUnix := rand.Int64N(endUnix-startUnix+1) + startUnix

	// Convert to time.Time
	return time.Unix(randomUnix, 0).UTC()
}
