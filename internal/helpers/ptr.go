package helpers

func Ptr[T any](v T) *T {
	return &v
}

func ResolvePointers[T any](list []*T) []T {
	result := make([]T, 0, len(list))
	for i := range list {
		result = append(result, ZeroIfNil(list[i]))
	}

	return result
}

func ZeroIfNil[T any](v *T) T {
	if v == nil {
		r := new(T)
		return *r
	}

	return *v
}
