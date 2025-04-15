package gengo

func Bool() bool {
	return UInt8Between(0, 1) == 1
}
