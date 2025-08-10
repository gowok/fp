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

func Test_Map(t *testing.T) {
	expected := []int{2, 4, 6}
	sample := []int{1, 2, 3}
	actual := Map(sample, func(s int) int {
		return s * 2
	})

	must := must.New(t)
	must.Equal(expected, actual)
}

func Test_MapIndex(t *testing.T) {
	expected := []int{2, 4, 6}
	sample := []int{1, 2, 3}
	actual := MapIndex(sample, func(s, i int) int {
		return s * 2
	})

	must := must.New(t)
	must.Equal(expected, actual)
}

func Test_Reduce(t *testing.T) {
	expected := 6
	sample := []int{1, 2, 3}
	actual := Reduce(sample, func(acc, val, i int) int {
		return acc + val
	}, 0)

	must := must.New(t)
	must.Equal(expected, actual)
}

func Test_Range(t *testing.T) {
	t.Run("max only", func(t *testing.T) {
		expected := []int{0, 1, 2, 3, 4}
		actual := Range(5)

		must := must.New(t)
		must.Equal(expected, actual)
	})

	t.Run("min and max", func(t *testing.T) {
		expected := []int{1, 2, 3, 4}
		actual := Range(1, 5)

		must := must.New(t)
		must.Equal(expected, actual)
	})

	t.Run("min, max, and steps", func(t *testing.T) {
		expected := []int{1, 3, 5, 7, 9}
		actual := Range(1, 10, 2)

		must := must.New(t)
		must.Equal(expected, actual)
	})

}

func Test_Contains(t *testing.T) {
	sample := []int{1, 2, 3}
	t.Run("positive", func(t *testing.T) {
		actual := Contains(sample, 1)

		must := must.New(t)
		must.True(actual)
	})

	t.Run("negative", func(t *testing.T) {
		actual := Contains(sample, 0)

		must := must.New(t)
		must.False(actual)
	})

}

func Test_Zip(t *testing.T) {
	sample1 := []int{1, 2, 3}
	sample2 := []int{1, 2, 3}

	must := must.New(t)
	for s1, s2 := range Zip(sample1, sample2) {
		must.Equal(s1, s2)
	}
}

func Test_Repeat(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := []byte("***")
		actual := Repeat('*', 3)

		must := must.New(t)
		must.Equal(len(expected), len(actual))
		for e, a := range Zip(expected, actual) {
			must.Equal(int(e), int(a))
		}
	})

	t.Run("negative", func(t *testing.T) {
		expected := []byte{}
		actual := Repeat('*', 0)

		must := must.New(t)
		must.Equal(len(expected), len(actual))

		actual = Repeat('*', -1)
		must.Equal(len(expected), len(actual))
	})
}

func Test_Uniq(t *testing.T) {
	t.Run("all different", func(t *testing.T) {
		sample := []int{1, 2, 3}
		output := Uniq(sample)

		must.New(t).Equal(sample, output)
	})

	t.Run("should remove some", func(t *testing.T) {
		sample := []int{1, 2, 2, 1, 3}
		output := Uniq(sample)
		expected := []int{1, 2, 3}

		must.New(t).Equal(expected, output)
	})
}
