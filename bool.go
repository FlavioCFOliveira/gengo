package gengo

import "math/rand/v2"

func Bool() bool {
	return rand.Uint32()&1 == 1
}
