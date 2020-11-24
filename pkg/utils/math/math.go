package math

func MinUint(x, y uint) uint {
	if x > y {
		return y
	}
	return x
}

func MaxUint(x, y uint) uint {
	if x < y {
		return y
	}
	return x
}
