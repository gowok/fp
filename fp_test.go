package fp

import (
	"testing"

	"github.com/golang-must/must"
)

func TestFlip(t *testing.T) {
	concat := func(a, b string) string { return a + b }
	must.Equal(t, "helloworld", concat("hello", "world"))
	must.Equal(t, "worldhello", Flip(concat)("hello", "world"))
}

func TestPipe(t *testing.T) {
	inc := func(a int) int { return a + 1 }
	double := func(a int) int { return a * 2 }
	must.Equal(t, 3, Pipe(3))
	must.Equal(t, 4, Pipe(3, inc))
	must.Equal(t, 8, Pipe(3, inc, double))
}
