package helpers

func Ternary[T any](flag bool, first T, second T) T {
	if flag {
		return first
	}

	return second
}
