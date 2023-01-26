package maps

func Keys[T comparable, U any](input map[T]U) []T {
	result := []T{}
	for k := range input {
		result = append(result, k)
	}

	return result
}

func Values[T comparable, U any](input map[T]U) []U {
	result := []U{}
	for _, v := range input {
		result = append(result, v)
	}

	return result
}
