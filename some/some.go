package some

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

type Some[T any] struct {
	value     *T
	isPresent bool
}

func newOptional[T any](val *T) Some[T] {
	isPresent := val != nil
	return Some[T]{val, isPresent}
}

func Empty[T any]() Some[T] {
	return newOptional[T](nil)
}

func Of[T any](val T) Some[T] {
	v := reflect.ValueOf(val)
	k := v.Kind()

	if k == reflect.Invalid {
		return Empty[T]()
	}

	if (k == reflect.Func || k == reflect.Pointer) && (v.IsZero() || v.IsNil()) {
		return Empty[T]()
	}

	return newOptional(&val)
}

func Map[T, U any](o Some[T], cb func(val T) U) Some[U] {
	if val, ok := o.Get(); ok {
		return Of(cb(val))
	}

	return Empty[U]()
}

func Filter[T any](o Some[T], cb func(val T) bool) Some[T] {
	if val, ok := o.Get(); ok && cb(val) {
		return o
	}

	return Empty[T]()
}

func FlatMap[T, U any](o Some[T], cb func(val T) Some[U]) Some[U] {
	if val, ok := o.Get(); ok {
		return cb(val)
	}

	return Empty[U]()
}

func (o Some[T]) Get() (T, bool) {
	if !o.IsPresent() {
		var value T
		return value, false
	}
	return *o.value, true
}

func (o Some[T]) IsPresent() bool {
	return o.isPresent
}

func (o Some[T]) IfPresent(callback func(T), elseCallback ...func()) {
	if val, ok := o.Get(); ok {
		callback(val)
		return
	}

	if len(elseCallback) > 0 {
		elseCallback[0]()
	}
}

func (o Some[T]) OrElse(val T) T {
	if !o.IsPresent() {
		return val
	}

	return *o.value
}

func (o Some[T]) OrElseFunc(gen func() T) T {
	if !o.IsPresent() {
		if gen != nil {
			return gen()
		}

		var zero T
		return zero
	}

	return *o.value
}

var ErrNoValuePresent = errors.New("no value present")

func (o Some[T]) OrPanic(errs ...error) T {
	if !o.IsPresent() {
		err := ErrNoValuePresent
		if len(errs) > 0 && errs[0] != nil {
			err = errs[0]
		}

		panic(err)
	}

	return *o.value
}

func (o *Some[T]) UnmarshalJSON(b []byte) error {
	o.value = nil
	o.isPresent = false

	if bytes.Equal(b, []byte("{}")) {
		return nil
	}

	var v *T
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v != nil {
		o.value = v
		o.isPresent = true
	}
	return nil
}

func (o Some[T]) MarshalJSON() ([]byte, error) {
	if o.IsPresent() {
		return json.Marshal(o.value)
	}

	return json.Marshal(map[string]any{})
}
