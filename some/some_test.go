package some

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang-must/must"
)

func TestEmpty(t *testing.T) {
	s := Empty[string]()
	must.False(t, s.IsPresent())
	val, ok := s.Get()
	must.False(t, ok)
	must.Equal(t, "", val)
}

func TestOf(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		isPresent bool
	}{
		{"positive/string", "hello", true},
		{"positive/int", 123, true},
		{"negative/nil func", (func())(nil), false},
		{"positive/valid func", func() {}, true},
		{"negative/nil pointer", (*int)(nil), false},
		{"positive/valid pointer", new(int), true},
		{"negative/nil interface", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Of(tt.input)
			must.Equal(t, tt.isPresent, s.IsPresent())
			if tt.isPresent {
				val, ok := s.Get()
				must.True(t, ok)
				must.Equal(t, tt.input, val)
			}
		})
	}
}

func TestSome_IsPresent(t *testing.T) {
	must.True(t, Of("test").IsPresent())
	must.False(t, Empty[string]().IsPresent())
}

func TestSome_IsEmpty(t *testing.T) {
	must.True(t, Empty[string]().IsEmpty())
	must.False(t, Of("test").IsEmpty())
}

func TestMap(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		actual := Map(Of(2), func(val int) int {
			return val * 10
		})

		must.True(t, actual.IsPresent())
		val, ok := actual.Get()
		must.True(t, ok)
		must.Equal(t, 20, val)
	})

	t.Run("empty", func(t *testing.T) {
		actual := Map(Empty[int](), func(val int) int {
			return val * 10
		})

		must.False(t, actual.IsPresent())
	})

	t.Run("change value type", func(t *testing.T) {
		val, ok := Map(Of(42), func(val int) string {
			return "answer"
		}).Get()

		must.True(t, ok)
		must.Equal(t, "answer", val)
	})
}

func TestFilter(t *testing.T) {
	t.Run("kept", func(t *testing.T) {
		actual := Filter(Of(10), func(val int) bool {
			return val > 5
		})

		must.True(t, actual.IsPresent())
		val, ok := actual.Get()
		must.True(t, ok)
		must.Equal(t, 10, val)
	})

	t.Run("dropped by predicate", func(t *testing.T) {
		actual := Filter(Of(2), func(val int) bool {
			return val > 5
		})

		must.False(t, actual.IsPresent())
	})

	t.Run("empty", func(t *testing.T) {
		actual := Filter(Empty[int](), func(val int) bool {
			return true
		})

		must.False(t, actual.IsPresent())
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("present returns mapped some", func(t *testing.T) {
		actual := FlatMap(Of(10), func(val int) Some[string] {
			return Of("ten")
		})

		val, ok := actual.Get()
		must.True(t, ok)
		must.Equal(t, "ten", val)
	})

	t.Run("present but mapping returns empty", func(t *testing.T) {
		actual := FlatMap(Of(10), func(val int) Some[string] {
			return Empty[string]()
		})

		must.False(t, actual.IsPresent())
	})

	t.Run("empty", func(t *testing.T) {
		actual := FlatMap(Empty[int](), func(val int) Some[string] {
			return Of("x")
		})

		must.False(t, actual.IsPresent())
	})
}

func TestSome_OrElse(t *testing.T) {
	must.Equal(t, "test", Of("test").OrElse("default"))
	must.Equal(t, "default", Empty[string]().OrElse("default"))
}

func TestSome_OrElseFunc(t *testing.T) {
	must.Equal(t, "test", Of("test").OrElseFunc(func() string { return "default" }))
	must.Equal(t, "default", Empty[string]().OrElseFunc(func() string { return "default" }))
	must.Equal(t, "", Empty[string]().OrElseFunc(nil))
}

func TestSome_Get(t *testing.T) {
	tests := []struct {
		name    string
		s       Some[string]
		wantVal string
		wantOk  bool
	}{
		{"positive/present", Of("test"), "test", true},
		{"negative/not present", Empty[string](), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.s.Get()
			must.Equal(t, tt.wantOk, ok)
			must.Equal(t, tt.wantVal, val)
		})
	}
}

func TestSome_IfPresent(t *testing.T) {
	tests := []struct {
		name            string
		s               Some[string]
		hasElseCallback bool
		wantCalled      bool
		wantElseCalled  bool
	}{
		{"positive/present with callback", Of("test"), true, true, false},
		{"negative/not present with else callback", Empty[string](), true, false, true},
		{"negative/not present without else callback", Empty[string](), false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			elseCalled := false

			callback := func(v string) {
				called = true
				must.Equal(t, "test", v)
			}

			if tt.hasElseCallback {
				tt.s.IfPresent(callback, func() {
					elseCalled = true
				})
			} else {
				tt.s.IfPresent(callback)
			}

			must.Equal(t, tt.wantCalled, called)
			must.Equal(t, tt.wantElseCalled, elseCalled)
		})
	}
}

func TestSome_OrPanic(t *testing.T) {
	tests := []struct {
		name      string
		s         Some[string]
		errs      []error
		wantVal   string
		wantPanic bool
		panicErr  string
	}{
		{
			name:    "positive/present",
			s:       Of("test"),
			wantVal: "test",
		},
		{
			name:      "negative/not present default panic",
			s:         Empty[string](),
			wantPanic: true,
		},
		{
			name:      "negative/not present custom panic",
			s:         Empty[string](),
			errs:      []error{errors.New("custom error")},
			wantPanic: true,
			panicErr:  "custom error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					r := recover()
					must.NotNil(t, r)
					if tt.panicErr != "" {
						must.Equal(t, tt.panicErr, r.(error).Error())
					}
				}()
			}

			val := tt.s.OrPanic(tt.errs...)
			if !tt.wantPanic {
				must.Equal(t, tt.wantVal, val)
			}
		})
	}
}

func TestSome_OrPanic_ErrorIdentity(t *testing.T) {
	t.Run("nil custom error on present value", func(t *testing.T) {
		must.Equal(t, "test", Of("test").OrPanic(nil))
	})

	t.Run("custom error is preserved", func(t *testing.T) {
		sentinel := errors.New("sentinel")
		defer func() {
			r := recover()
			must.True(t, errors.Is(r.(error), sentinel))
		}()
		Empty[string]().OrPanic(sentinel)
	})
}

func TestSome_JSON(t *testing.T) {
	type container struct {
		Opt Some[string] `json:"opt"`
	}

	tests := []struct {
		name      string
		operation string
		input     any
		want      string
		wantErr   bool
	}{
		{
			name:      "positive/marshal present",
			operation: "marshal",
			input:     Of("test"),
			want:      "\"test\"",
		},
		{
			name:      "positive/marshal empty",
			operation: "marshal",
			input:     Empty[string](),
			want:      "{}",
		},
		{
			name:      "positive/marshal embedded",
			operation: "marshal",
			input:     container{Opt: Of("test")},
			want:      "{\"opt\":\"test\"}",
		},
		{
			name:      "positive/unmarshal present",
			operation: "unmarshal",
			input:     "\"test\"",
			want:      "test",
		},
		{
			name:      "positive/unmarshal empty object",
			operation: "unmarshal",
			input:     "{}",
		},
		{
			name:      "positive/unmarshal null",
			operation: "unmarshal",
			input:     "null",
		},
		{
			name:      "negative/unmarshal error",
			operation: "unmarshal",
			input:     "{invalid}",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.operation == "marshal" {
				b, err := json.Marshal(tt.input)
				must.Nil(t, err)
				must.Equal(t, tt.want, string(b))
			} else {
				var s Some[string]
				err := json.Unmarshal([]byte(tt.input.(string)), &s)
				if tt.wantErr {
					must.NotNil(t, err)
				} else {
					must.Nil(t, err)
					if tt.want != "" {
						must.True(t, s.IsPresent())
						val, _ := s.Get()
						must.Equal(t, tt.want, val)
					} else {
						must.False(t, s.IsPresent())
					}
				}
			}
		})
	}

	t.Run("positive/unmarshal null overwrites present value", func(t *testing.T) {
		s := Of("test")
		err := json.Unmarshal([]byte("null"), &s)
		must.Nil(t, err)
		must.False(t, s.IsPresent())
	})
}
