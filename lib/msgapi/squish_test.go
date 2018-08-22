package msgapi

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestSquishBufHash32(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check Squish bufHash32()", func() {
		g.It("calculates correct Sqush bufhash-32 of an empty string", func() {
			g.Assert(bufHash32("")).Equal(uint32(0))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Alexander N. Skovpen'", func() {
			g.Assert(bufHash32("Alexander N. Skovpen")).Equal(uint32(0x00efd7be))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Эдуардыч'[CP866]", func() {
			g.Assert(bufHash32(string([]byte{'\x9D', '\xA4', '\xE3', '\xA0', '\xE0', '\xA4', '\xEB', '\xE7'}))).Equal(uint32(0x02debb97))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Юрий Григорьев'[CP866]", func() {
			g.Assert(bufHash32(string([]byte{'\x9E', '\xE0', '\xA8', '\xA9', '\x20', '\x83', '\xE0', '\xA8', '\xA3', '\xAE', '\xE0', '\xEC', '\xA5', '\xA2'}))).Equal(uint32(0x7c100ff2))
		})
	})
}
