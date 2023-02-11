package slices

import (
	"testing"

	"github.com/golang-must/must"
)

func Test_Foreach(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		sample := []int{1, 2, 3}
		must := must.New(t)

		ForEach(sample, func(s, i int) {
			must.Equal(sample[i], s)
		})
	})
}

func Test_GoForeach(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		sample := []int{1, 2, 3}
		must := must.New(t)

		GoForEach(sample, func(s, i int) {
			must.Equal(sample[i], s)
		})
	})
}
