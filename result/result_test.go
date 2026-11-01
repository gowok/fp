package result

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang-must/must"
)

func TestOkUnwrap(t *testing.T) {
	r := Ok(123)
	must.True(t, r.IsOk())
	v, err := r.Unwrap()
	must.Nil(t, err)
	must.Equal(t, 123, v)
}

func TestErrUnwrap(t *testing.T) {
	e := errors.New("boom")
	r := Err[int](e)
	must.True(t, r.IsErr())
	_, err := r.Unwrap()
	must.NotNil(t, err)
	must.Equal(t, "boom", err.Error())
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name    string
		create  func() Result[string]
		wantVal string
		wantErr bool
	}{
		{"ok", func() Result[string] { return Ok("x") }, "x", false},
		{"err", func() Result[string] { return Err[string](errors.New("nope")) }, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.create()
			if tt.wantErr {
				must.True(t, r.IsErr())
			} else {
				must.True(t, r.IsOk())
				v, err := r.Unwrap()
				must.Nil(t, err)
				must.Equal(t, tt.wantVal, v)
			}
		})
	}
}

func TestUnwrapOr(t *testing.T) {
	r := Err[string](errors.New("e"))
	got := r.UnwrapOr("def")
	must.Equal(t, "def", got)
}

func TestMap(t *testing.T) {
	r := Ok(2)
	r2 := Map(r, func(x int) string { return fmt.Sprintf("%d", x*3) })
	v, err := r2.Unwrap()
	must.Nil(t, err)
	must.Equal(t, "6", v)

	re := Err[int](errors.New("bad"))
	r3 := Map(re, func(x int) string { return "x" })
	must.True(t, r3.IsErr())
}

func TestErrMethod(t *testing.T) {
	e := errors.New("boom")
	r := Err[int](e)
	must.Equal(t, e, r.Err())
	ok := Ok(5)
	must.Nil(t, ok.Err())
}

func TestMapAbsentResult(t *testing.T) {
	r := New[string](nil, nil)
	mapped := Map(r, func(string) string { return "unexpected" })

	must.True(t, mapped.IsErr())
	must.Equal(t, ErrNoValuePresent, mapped.Err())
}

func TestUnwrapOr_PresentAndAbsent(t *testing.T) {
	must.Equal(t, "val", Ok("val").UnwrapOr("def"))
	must.Equal(t, 10, Ok(10).UnwrapOr(0))

	must.Equal(t, "def", New[string](nil, nil).UnwrapOr("def"))
	must.Equal(t, "def", Err[string](errors.New("e")).UnwrapOr("def"))
}

func TestOrPanic(t *testing.T) {
	t.Run("ok returns value", func(t *testing.T) {
		must.Equal(t, 7, Ok(7).OrPanic())
	})

	t.Run("err panics with stored error", func(t *testing.T) {
		e := errors.New("boom")
		defer func() {
			must.Equal(t, e, recover())
		}()
		Err[int](e).OrPanic()
	})

	t.Run("absent panics with no value present", func(t *testing.T) {
		defer func() {
			must.Equal(t, ErrNoValuePresent, recover())
		}()
		New[int](nil, nil).OrPanic()
	})

	t.Run("custom error overrides stored error", func(t *testing.T) {
		custom := errors.New("custom")
		defer func() {
			must.Equal(t, custom, recover())
		}()
		Err[int](errors.New("boom")).OrPanic(custom)
	})

	t.Run("nil custom error falls back to stored error", func(t *testing.T) {
		e := errors.New("boom")
		defer func() {
			must.Equal(t, e, recover())
		}()
		Err[int](e).OrPanic(nil)
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("ok to ok", func(t *testing.T) {
		r := FlatMap(Ok(2), func(x int) Result[string] {
			return Ok(fmt.Sprintf("%d", x*3))
		})

		v, err := r.Unwrap()
		must.Nil(t, err)
		must.Equal(t, "6", v)
	})

	t.Run("ok to err", func(t *testing.T) {
		r := FlatMap(Ok(2), func(int) Result[string] {
			return Err[string](errors.New("bad"))
		})

		must.True(t, r.IsErr())
		must.Equal(t, "bad", r.Err().Error())
	})

	t.Run("err keeps original error", func(t *testing.T) {
		src := errors.New("boom")
		r := FlatMap(Err[int](src), func(int) Result[string] {
			return Ok("unexpected")
		})

		must.True(t, r.IsErr())
		must.Equal(t, src, r.Err())
	})

	t.Run("absent becomes no value present", func(t *testing.T) {
		r := FlatMap(New[string](nil, nil), func(string) Result[int] {
			return Ok(0)
		})

		must.True(t, r.IsErr())
		must.Equal(t, ErrNoValuePresent, r.Err())
	})
}

func TestOr(t *testing.T) {
	t.Run("ok keeps itself", func(t *testing.T) {
		r := Or(Ok(1), Ok(2))
		must.Equal(t, 1, r.UnwrapOr(0))
	})

	t.Run("err falls back", func(t *testing.T) {
		r := Or(Err[int](errors.New("boom")), Ok(2))
		must.Equal(t, 2, r.UnwrapOr(0))
	})

	t.Run("absent falls back", func(t *testing.T) {
		r := Or(New[int](nil, nil), Ok(3))
		must.Equal(t, 3, r.UnwrapOr(0))
	})
}

func TestOrElse(t *testing.T) {
	t.Run("ok skips callback", func(t *testing.T) {
		called := false
		r := OrElse(Ok(1), func() Result[int] {
			called = true
			return Ok(2)
		})

		must.False(t, called)
		must.Equal(t, 1, r.UnwrapOr(0))
	})

	t.Run("err calls callback", func(t *testing.T) {
		r := OrElse(Err[int](errors.New("boom")), func() Result[int] {
			return Ok(2)
		})

		must.True(t, r.IsOk())
		must.Equal(t, 2, r.UnwrapOr(0))
	})

	t.Run("absent calls callback", func(t *testing.T) {
		r := OrElse(New[int](nil, nil), func() Result[int] {
			return Ok(3)
		})

		must.True(t, r.IsOk())
		must.Equal(t, 3, r.UnwrapOr(0))
	})
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input Result[any]
		want  string
	}{
		{"present int", Ok[any](7), "7"},
		{"present string", Ok[any]("x"), `"x"`},
		{"error", Err[any](errors.New("boom")), `{"error":"boom"}`},
		{"absent", New[any](nil, nil), "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.input)
			must.Nil(t, err)
			must.Equal(t, tt.want, string(b))
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("present value", func(t *testing.T) {
		var r Result[int]
		must.Nil(t, json.Unmarshal([]byte("7"), &r))
		must.True(t, r.IsOk())
		must.Equal(t, 7, r.UnwrapOr(0))
	})

	t.Run("error payload", func(t *testing.T) {
		var r Result[int]
		must.Nil(t, json.Unmarshal([]byte(`{"error":"boom"}`), &r))
		must.True(t, r.IsErr())
		must.Equal(t, "boom", r.Err().Error())
	})

	t.Run("empty object is absent", func(t *testing.T) {
		var r Result[int]
		must.Nil(t, json.Unmarshal([]byte("{}"), &r))
		must.False(t, r.IsOk())
		must.False(t, r.IsErr())
		must.Equal(t, 0, r.UnwrapOr(0))
	})

	t.Run("null is absent", func(t *testing.T) {
		r := Ok(3)
		must.Nil(t, json.Unmarshal([]byte("null"), &r))
		must.False(t, r.IsOk())
	})

	t.Run("round trip value", func(t *testing.T) {
		b, err := json.Marshal(Ok(42))
		must.Nil(t, err)

		var r Result[int]
		must.Nil(t, json.Unmarshal(b, &r))
		must.Equal(t, 42, r.UnwrapOr(0))
	})

	t.Run("round trip error", func(t *testing.T) {
		b, err := json.Marshal(Err[string](errors.New("bad")))
		must.Nil(t, err)

		var r Result[string]
		must.Nil(t, json.Unmarshal(b, &r))
		must.True(t, r.IsErr())
		must.Equal(t, "bad", r.Err().Error())
	})
}
