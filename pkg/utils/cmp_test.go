package utils

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestNamesEqual(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check NamesEqual()", func() {
		g.It("check normal equals", func() {
			g.Assert(NamesEqual("Alexander Skovpen", " Alexander Skovpen ")).Equal(true)
		})
		g.It("check dot equals", func() {
			g.Assert(NamesEqual("Alexander N Skovpen", " Alexander N. Skovpen ")).Equal(true)
		})
		g.It("check false equals", func() {
			g.Assert(NamesEqual("Alexander N. Skovpen", " Alexander P. Skovpen ")).Equal(false)
		})
	})
}
