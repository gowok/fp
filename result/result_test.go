package result

import (
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

func TestUnwrapOr_PresentAndAbsent(t *testing.T) {
	must.Equal(t, "val", Ok("val").UnwrapOr("def"))
	must.Equal(t, 10, Ok(10).UnwrapOr(0))

	must.Equal(t, "def", New[string](nil, nil).UnwrapOr("def"))
	must.Equal(t, "def", Err[string](errors.New("e")).UnwrapOr("def"))
}
