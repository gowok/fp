package slices

import (
	"slices"
	"sync"
)

func ForEach[T any](input []T, cb func(val T, index int)) {
	for i, v := range input {
		cb(v, i)
	}
}

func GoForEach[T any](input []T, cb func(val T, index int)) {
	wg := sync.WaitGroup{}
	for i, v := range input {
		wg.Add(1)
		go func(vv T, j int) {
			defer wg.Done()
			cb(vv, j)
		}(v, i)
	}
	wg.Wait()
}

func Filter[T comparable](input []T, cb func(val T, index int) bool) []T {
	result := []T{}
	for i, v := range input {
		if cb(v, i) {
			result = append(result, v)
		}
	}

	return result
}

func Map[T comparable, U any](input []T, cb func(val T) U) []U {
	return MapIndex(input, func(val T, index int) U {
		return cb(val)
	})
}

func MapIndex[T comparable, U any](input []T, cb func(val T, index int) U) []U {
	result := []U{}
	for i, v := range input {
		result = append(result, cb(v, i))
	}

	return result
}

func UniqMap[T comparable, U comparable](input []T, cb func(val T) U) []U {
	maps := Map(input, cb)
	return Uniq(maps)
}

func Reduce[T comparable, U any](input []T, cb func(acc U, val T, index int) U, initial U) U {
	var result U
	for i, v := range input {
		result = cb(result, v, i)
	}

	return result
}

func Range[T int | float32 | float64](input T, params ...T) []T {
	result := []T{}
	var start T = 0
	var finish T = input
	var step T = 1

	if len(params) > 0 {
		start = input
		finish = params[0]
	}

	if len(params) == 2 {
		step = params[1]
	}

	for ; start < finish; start += step {
		result = append(result, start)
	}

	return result
}

func Contains[T comparable](input []T, comp T) bool {
	return slices.Contains(input, comp)
}

func Zip[T any, U any](slice1 []T, slice2 []U) func(yield func(T, U) bool) {
	return func(yield func(x T, y U) bool) {
		for i, ii := range slice1 {
			res := yield(ii, slice2[i])
			if !res {
				break
			}
		}
	}
}

func Repeat[T any](input T, times int) (output []T) {
	if times <= 0 {
		return
	}

	for range times {
		output = append(output, input)
	}

	return
}

func Uniq[T comparable](input []T) []T {
	output := make([]T, 0, len(input))
	track := make(map[T]struct{}, len(input))
	for i := range input {
		if _, ok := track[input[i]]; ok {
			continue
		}

		track[input[i]] = struct{}{}
		output = append(output, input[i])
	}

	return output
}

func Compact[T comparable](input []T) []T {
	var zero T
	output := make([]T, 0, len(input))

	for i := range input {
		if input[i] != zero {
			output = append(output, input[i])
		}
	}

	return output
}
