package gengo

import (
	"math/rand"
	"time"
)

func Date() time.Time {
	randomUnix := rand.Int63n(dtMax-dtMin+1) + dtMin
	return time.Unix(randomUnix, 0).UTC()
}
func UnixDate() time.Time {
	randomUnix := rand.Int63n(dtUnixMax-dtUnixMin+1) + dtUnixMin
	return time.Unix(randomUnix, 0).UTC()
}

func DateBetween(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}

	startUnix := start.Unix()
	endUnix := end.Unix()

	// Garante que o range não está vazio
	if startUnix == endUnix {
		return start
	}

	// Gera um número aleatório entre startUnix e endUnix
	randomUnix := rand.Int63n(endUnix-startUnix+1) + startUnix

	// Converte para time.Time
	return time.Unix(randomUnix, 0).UTC()
}
