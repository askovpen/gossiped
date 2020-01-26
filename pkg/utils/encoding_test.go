package utils

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestCharmapEqual(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check DecodeCharmap()", func() {
		g.It("check cp866", func() {
			g.Assert(DecodeCharmap("\x92\xa5\xe1\xe2", "CP866")).Equal("Тест")
		})
		g.It("check utf-8", func() {
			g.Assert(DecodeCharmap("Тест", "UTF-8")).Equal("Тест")
		})
	})
	g.Describe("Check EncodeCharmap()", func() {
		g.It("check cp866", func() {
			g.Assert(EncodeCharmap("Тест", "CP866")).Equal("\x92\xa5\xe1\xe2")
		})
		g.It("check utf-8", func() {
			g.Assert(EncodeCharmap("Тест", "UTF-8")).Equal("Тест")
		})
	})
}
