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

func Map[T comparable, U any](input []T, cb func(val T, index int) U) []U {
	result := []U{}
	for i, v := range input {
		result = append(result, cb(v, i))
	}

	return result
}

func Reduce[T comparable, U any](input []T, cb func(acc U, val T, index int) U, initial U) U {
	var result U
	for i, v := range input {
		result = cb(result, v, i)
	}

	return result
}
