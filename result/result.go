package result

import (
	"bytes"
	"encoding/json"
	"errors"
)

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

func (r Result[T]) OrPanic(errs ...error) T {
	v, err := r.Unwrap()
	if err != nil {
		if len(errs) > 0 && errs[0] != nil {
			err = errs[0]
		}
		panic(err)
	}
	return v
}

func Map[T any, U any](r Result[T], fn func(T) U) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	if !r.isPresent {
		return Err[U](ErrNoValuePresent)
	}
	return Ok(fn(*r.value))
}

func FlatMap[T any, U any](r Result[T], fn func(T) Result[U]) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	if !r.isPresent {
		return Err[U](ErrNoValuePresent)
	}
	return fn(*r.value)
}

func Or[T any](r Result[T], alt Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return alt
}

func OrElse[T any](r Result[T], fn func() Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return fn()
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.err != nil {
		return json.Marshal(map[string]any{"error": r.err.Error()})
	}
	if !r.isPresent {
		return []byte("null"), nil
	}
	return json.Marshal(r.value)
}

func (r *Result[T]) UnmarshalJSON(b []byte) error {
	r.value = nil
	r.err = nil
	r.isPresent = false

	if bytes.Equal(b, []byte("{}")) || bytes.Equal(b, []byte("null")) {
		return nil
	}

	var envelope struct {
		Err *string `json:"error"`
	}
	if err := json.Unmarshal(b, &envelope); err == nil && envelope.Err != nil {
		r.err = errors.New(*envelope.Err)
		return nil
	}

	var v *T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v != nil {
		r.value = v
		r.isPresent = true
	}
	return nil
}
