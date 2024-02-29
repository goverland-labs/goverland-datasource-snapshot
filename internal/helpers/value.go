package helpers

func ValurOrDefault[T any](v *T, d T) T {
	if v == nil {
		return d
	}

	return *v
}
