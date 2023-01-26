package slices

func Filter[T comparable](input []T, cb func(val T, index int) bool) []T {
	result := []T{}
	for i, v := range input {
		if cb(v, i) {
			result = append(result, v)
		}
	}

	return result
}
