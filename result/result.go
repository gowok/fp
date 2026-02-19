package result

import "errors"

type Result[T any] struct {
	value     *T
	err       error
	isPresent bool
}

func New[T any](val *T, err error) Result[T] {
	isPresent := val != nil && err == nil
	return Result[T]{value: val, err: err, isPresent: isPresent}
}

func Ok[T any](v T) Result[T] { return New(&v, nil) }

func Err[T any](err error) Result[T] { return New[T](nil, err) }

func (r Result[T]) IsOk() bool { return r.err == nil && r.isPresent }

func (r Result[T]) IsErr() bool { return r.err != nil }

func (r Result[T]) Err() error { return r.err }

var ErrNoValuePresent = errors.New("no value present")

func (r Result[T]) Unwrap() (T, error) {
	var zero T
	if r.err != nil {
		return zero, r.err
	}
	if !r.isPresent {
		return zero, ErrNoValuePresent
	}
	return *r.value, nil
}

func (r Result[T]) UnwrapOr(def T) T {
	if r.err != nil || !r.isPresent {
		return def
	}
	return *r.value
}

func Map[T any, U any](r Result[T], fn func(T) U) Result[U] {
	if r.err != nil || !r.isPresent {
		return Err[U](r.err)
	}
	return Ok(fn(*r.value))
}
