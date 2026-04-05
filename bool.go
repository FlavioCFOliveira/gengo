package gengo

import "math/rand/v2"

// Bool returns a random boolean value.
func Bool() bool {
	return rand.Uint32()&1 == 1
}
