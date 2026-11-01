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
