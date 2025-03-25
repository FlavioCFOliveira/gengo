package gengo

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestReseed(t *testing.T) {
	Reseed()
}

func TestReseedFunc(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		f := func() *rand.Rand {
			return rand.New(rand.NewSource(time.Now().UnixNano()))
		}

		if err := ReseedFunc(f); err != nil {
			t.Errorf("ReseedFunc() error = %v", err)
		}

	})

	t.Run("NilFunc", func(t *testing.T) {

		if err := ReseedFunc(nil); !errors.Is(err, ErrorNilFunc) {
			t.Errorf("ReseedFunc() error = %v expected = %v", err, ErrorNilFunc)
		}

	})

	t.Run("FuncNilResult", func(t *testing.T) {
		f := func() *rand.Rand {
			return nil
		}

		if err := ReseedFunc(f); !errors.Is(err, ErrorFuncNilResult) {
			t.Errorf("ReseedFunc() error = %v expected = %v", err, ErrorFuncNilResult)
		}

	})

}
