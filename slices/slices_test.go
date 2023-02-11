package slices

import (
	"testing"

	"github.com/golang-must/must"
)

func Test_Foreach(t *testing.T) {
	sample := []int{1, 2, 3}
	must := must.New(t)

	ForEach(sample, func(s, i int) {
		must.Equal(sample[i], s)
	})
}

func Test_GoForeach(t *testing.T) {
	sample := []int{1, 2, 3}
	must := must.New(t)

	GoForEach(sample, func(s, i int) {
		must.Equal(sample[i], s)
	})
}

func Test_Filter(t *testing.T) {
	expected := []int{1}
	sample := []int{1, 2, 3}
	actual := Filter(sample, func(s, i int) bool {
		return s == 1
	})

	must := must.New(t)
	must.Equal(expected, actual)
}
