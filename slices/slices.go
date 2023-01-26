package slices

import "sync"

func ForEach[T any](input []T, cb func(val T, index int)) {
	for i, v := range input {
		cb(v, i)
	}
}

func GoForEach[T any](input []T, cb func(val T, index int)) {
	wg := sync.WaitGroup{}
	for i, v := range input {
		wg.Add(1)
		go func(vv T) {
			defer wg.Done()
			cb(vv, i)
		}(v)
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

func Includes[T comparable](input []T, comp T) bool {
	for _, v := range input {
		if comp == v {
			return true
		}
	}

	return false
}
