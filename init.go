package gengo

import (
	"errors"
	"time"
)

var (
	ErrorNilFunc       = errors.New("f cannot be nil")
	ErrorFuncNilResult = errors.New("f returned nil")
)

var dtMin, dtMax, dtUnixMin, dtUnixMax int64

func init() {
	dtMin = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	dtMax = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	dtUnixMin = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	dtUnixMax = time.Date(2038, 01, 19, 3, 14, 07, 0, time.UTC).Unix()
}
