package maps

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/golang-must/must"
)

var sample = map[string]int{
	"x": 1,
	"y": 2,
	"z": 3,
}
var keys = []string{"x", "y", "z"}
var values = []int{1, 2, 3}

func Test_Keys(t *testing.T) {
	actual := Keys(sample)

	must := must.New(t)
	must.Equal(len(keys), len(actual))
}

func Test_Values(t *testing.T) {
	actual := Values(sample)

	must := must.New(t)
	must.Equal(len(values), len(actual))
}

func Test_KeysSeq(t *testing.T) {
	must := must.New(t)
	seen := map[string]bool{}
	for k := range KeysSeq(sample) {
		seen[k] = true
	}

	must.Equal(len(sample), len(seen))
	must.True(seen["x"])
	must.True(seen["y"])
	must.True(seen["z"])
}

func Test_ValuesSeq(t *testing.T) {
	must := must.New(t)
	count := 0
	for range ValuesSeq(sample) {
		count++
	}
	must.Equal(len(sample), count)
}

func Test_Has(t *testing.T) {
	must := must.New(t)
	must.True(Has(sample, "x"))
	must.False(Has(sample, "nope"))
	must.False(Has(map[string]int{}, "x"))
}

func Test_IsEmpty(t *testing.T) {
	must := must.New(t)
	must.True(IsEmpty(map[string]int{}))
	must.True(IsEmpty(map[string]int(nil)))
	must.False(IsEmpty(sample))
}

func Test_Clone(t *testing.T) {
	must := must.New(t)

	clone := Clone(sample)
	must.Equal(sample, clone)

	clone["x"] = 100
	must.Equal(1, sample["x"])

	must.Equal(map[string]int{}, Clone(map[string]int(nil)))
}

func Test_CopyBy(t *testing.T) {
	expected := map[string]int{
		"x": 1,
		"y": 2,
	}
	actual := CopyBy(sample, func(key string, value int) bool {
		return key != "z"
	})

	must := must.New(t)
	must.Equal(expected, actual)
}

func Test_Filter(t *testing.T) {
	expected := map[string]int{
		"y": 2,
		"z": 3,
	}
	actual := Filter(sample, func(value int) bool {
		return value > 1
	})

	must := must.New(t)
	must.Equal(expected, actual)
	must.Equal(map[string]int{"x": 1, "y": 2, "z": 3}, sample)
}

func Test_Map(t *testing.T) {
	expected := map[string]string{
		"x": "1!",
		"y": "2!",
		"z": "3!",
	}
	actual := Map(sample, func(value int) string {
		return strconv.Itoa(value) + "!"
	})

	must := must.New(t)
	must.Equal(expected, actual)
	must.Equal(map[string]int{"x": 1, "y": 2, "z": 3}, sample)
}

func Test_MapEmpty(t *testing.T) {
	actual := Map[string, int, string](nil, func(value int) string {
		return strconv.Itoa(value)
	})

	must := must.New(t)
	must.Equal(map[string]string{}, actual)
	must.NotNil(actual)
}

func Test_CopyByKeys(t *testing.T) {
	expected := map[string]int{
		"x": 1,
		"y": 2,
	}
	actual := CopyByKeys(sample, []string{"x", "y"})

	must := must.New(t)
	must.Equal(expected, actual)
}

func Test_Omit(t *testing.T) {
	expected := map[string]int{
		"x": 1,
		"y": 2,
	}
	actual := Omit(sample, []string{"z"})

	m := must.New(t)
	m.Equal(expected, actual)
	m.Equal(map[string]int{"x": 1, "y": 2, "z": 3}, sample)
}

func Test_Entries(t *testing.T) {
	expected := []Entry[string, int]{
		{"x", 1},
		{"y", 2},
		{"z", 3},
	}
	actual := Entries(sample)

	must := must.New(t)
	must.Equal(len(expected), len(actual))
}

func Test_EntriesSeq(t *testing.T) {
	must := must.New(t)
	count := 0
	for range EntriesSeq(sample) {
		count++
	}
	must.Equal(len(sample), count)
}

func Test_FromEntries(t *testing.T) {
	entries := []Entry[string, int]{
		{"x", 1},
		{"y", 2},
		{"z", 3},
	}

	must := must.New(t)
	must.Equal(sample, FromEntries(entries))
	must.Equal(map[string]int{}, FromEntries[string, int](nil))
}

func Test_Combine(t *testing.T) {
	expected := map[string]int{
		"x": 1,
		"y": 2,
		"z": 3,
		"a": 4,
		"b": 5,
	}
	actual := Combine(sample, map[string]int{
		"a": 4,
		"b": 5,
	})

	must := must.New(t)
	must.Equal(expected, actual)
	must.Equal(map[string]int{"x": 1, "y": 2, "z": 3}, sample)
}

func Test_FromStruct(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := map[string]float64{
			"x": 1,
			"y": 2,
			"z": 3,
			"a": 4,
			"b": 5,
		}
		actual := FromStruct(struct {
			A int `json:"a"`
			B int `json:"b"`
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		}{4, 5, 1, 2, 3})

		must := must.New(t)
		for k, v := range actual {
			ev, ok := expected[k]
			must.True(ok)
			must.Equal(ev, v)
		}
	})

	t.Run("primitive", func(t *testing.T) {
		actual := FromStruct(1)
		must.New(t).Equal("map[string]interface {}", reflect.TypeOf(actual).String())
	})

	t.Run("func", func(t *testing.T) {
		actual := FromStruct(func() {})
		must.New(t).Equal("map[string]interface {}", reflect.TypeOf(actual).String())
	})
}

func Test_ToStruct(t *testing.T) {
	type example struct {
		A int `json:"a"`
		B int `json:"b"`
		X int `json:"x"`
		Y int `json:"y"`
		Z int `json:"z"`
	}
	t.Run("positive", func(t *testing.T) {
		input := map[string]any{
			"x": 1,
			"y": 2,
			"z": 3,
			"a": 4,
			"b": 5,
		}
		actual := example{}
		expected := example{4, 5, 1, 2, 3}
		err := ToStruct(input, &actual)

		must := must.New(t)
		must.Nil(err)
		must.Equal(actual, expected)
	})

	t.Run("func", func(t *testing.T) {
		err := ToStruct(map[string]any{"x": func() {}}, struct{}{})
		must.New(t).NotNil(err)
	})
}

func Test_UniqValues(t *testing.T) {
	input := map[string]int{
		"x": 1,
		"y": 2,
		"z": 1,
		"a": 3,
		"b": 2,
	}
	actual := UniqValues(input)
	must.New(t).Equal(len(actual), 3)
}

func Test_Get(t *testing.T) {
	type Map map[string]any
	must := must.New(t)
	var cases = []struct {
		Title string
		Input map[string]any
		Check func(*testing.T, map[string]any)
	}{
		{
			"flat",
			map[string]any{"host": "localhost", "enabled": true},
			func(t *testing.T, c map[string]any) {
				must.Equal(Get(c, "host", ""), "localhost")
				must.Equal(Get(c, "enabled", false), true)
			},
		},
		{
			"nested",
			map[string]any{"host": "localhost", "address": map[string]any{"city": "Jakarta"}},
			func(t *testing.T, c map[string]any) {
				must.Equal(Get(c, "host", ""), "localhost")
				must.Equal(Get(c, "address.city", ""), "Jakarta")
				must.Equal(Get(c, "host.name", ""), "")
			},
		},
		{
			"nested type alias",
			map[string]any{"host": "localhost", "address": Map{"city": "Jakarta"}},
			func(t *testing.T, c map[string]any) {
				must.Equal(Get(c, "host", ""), "localhost")
				must.Equal(Get(c, "address.city", ""), "Jakarta")
				must.Equal(Get(c, "host.name", ""), "")
			},
		},
		{
			"nested type any",
			map[string]any{"host": "localhost", "address": any(map[string]any{"city": "Jakarta"})},
			func(t *testing.T, c map[string]any) {
				must.Equal(Get(c, "host", ""), "localhost")
				must.Equal(Get(c, "address.city", ""), "Jakarta")
				must.Equal(Get(c, "host.name", ""), "")
			},
		},
		{
			"default if not possible",
			map[string]any{"host": "localhost", "enabled": func() bool { return true }},
			func(t *testing.T, c map[string]any) {
				must.Equal(Get(c, "host", ""), "localhost")
				must.Equal(Get(c, "enabled", false), false)
				must.Equal(Get(c, "uwu", 0), 0)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			c.Check(t, c.Input)
		})
	}
}
