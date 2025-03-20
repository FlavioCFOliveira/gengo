package GoGen

import (
	"errors"
	"math/rand"
	"time"
)

var rnd *rand.Rand
var ErrorNilFunc = errors.New("f cannot be nil")
var ErrorFuncNilResult = errors.New("f returned nil")

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Reseed() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func ReseedFunc(f func() *rand.Rand) error {
	if f == nil {
		return ErrorNilFunc
	}

	newRnd := f()
	if newRnd == nil {
		return ErrorFuncNilResult
	}
	rnd = newRnd

	return nil
}
