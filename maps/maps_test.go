package maps

import (
	"reflect"
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

func Test_ValuesSeq(t *testing.T) {
	must := must.New(t)
	count := 0
	for range ValuesSeq(sample) {
		count++
	}
	must.Equal(len(sample), count)
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

func Test_CopyByKeys(t *testing.T) {
	expected := map[string]int{
		"x": 1,
		"y": 2,
	}
	actual := CopyByKeys(sample, []string{"x", "y"})

	must := must.New(t)
	must.Equal(expected, actual)
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
