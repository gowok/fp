package maps

func Keys[T comparable, U any](input map[T]U) []T {
	result := []T{}
	for k, _ := range input {
		result = append(result, k)
	}

	return result
}
