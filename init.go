// Package gengo provides a simple and practical way to generate random data
// structures. It offers a zero-dependency API for generating random values of
// standard Go types including integers, floats, complex numbers, strings,
// booleans, dates, and words.
package gengo

import "time"

var dtMin, dtMax, dtUnixMin, dtUnixMax int64
var dtRange, dtUnixRange int64

func init() {
	dtMin = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	dtMax = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	dtUnixMin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	dtUnixMax = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC).Unix()
	dtRange = dtMax - dtMin + 1
	dtUnixRange = dtUnixMax - dtUnixMin + 1
}
