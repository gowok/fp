package maps

import (
	"encoding/json"
	"maps"

	"github.com/gowok/fp/slices"
)

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

func ValuesSeq[T comparable, U any](input map[T]U) func(yield func(v U) bool) {
	return func(yield func(v U) bool) {
		for _, v := range input {
			res := yield(v)
			if !res {
				break
			}
		}
	}
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
		return slices.Contains(keys, k)
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

func EntriesSeq[T comparable, U any](input map[T]U) func(yield func(e Entry[T, U]) bool) {
	return func(yield func(e Entry[T, U]) bool) {
		for k, v := range input {
			res := yield(Entry[T, U]{k, v})
			if !res {
				break
			}
		}
	}
}

func Combine[T comparable, U any](input1 map[T]U, input2 map[T]U) map[T]U {
	result := input1
	maps.Copy(result, input2)
	return result
}

func FromStruct[T any](input T) (output map[string]any) {
	jsonb, err := json.Marshal(input)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonb, &output)
	if err != nil {
		return
	}

	return
}
