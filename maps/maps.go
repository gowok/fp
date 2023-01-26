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

func CopyBy[T comparable, U any](input map[T]U, cb func(key T, value U) bool) map[T]U {
	result := map[T]U{}
	for k, v := range input {
		if cb(k, v) {
			result[k] = v
		}
	}

	return result
}

func CopyByKeys[T comparable, U any](input map[T]U, keys []T) map[T]U {
	return CopyBy(input, func(k T, v U) bool {
		for _, kk := range keys {
			if k == kk {
				return true
			}
		}

		return false
	})
}

type Entry[T comparable, U any] struct {
	Key   T
	Value U
}

func Entries[T comparable, U any](input map[T]U) []Entry[T, U] {
	result := []Entry[T, U]{}
	for k, v := range input {
		result = append(result, Entry[T, U]{k, v})
	}

	return result
}
