package maps

import (
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
