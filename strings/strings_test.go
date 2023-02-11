package strings

import (
	"testing"

	"github.com/golang-must/must"
)

func Test_Repeat(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := "***"
		actual := Repeat("*", 3)

		must := must.New(t)
		must.Equal(expected, actual)
	})

	t.Run("negative", func(t *testing.T) {
		expected := ""
		actual := Repeat("*", 0)

		must := must.New(t)
		must.Equal(expected, actual)

		actual = Repeat("*", -1)
		must.Equal(expected, actual)
	})
}
