package msgapi

import(
  "testing"
  . "github.com/franela/goblin"
)

func TestSquishBufHash32(t *testing.T) {
  g := Goblin(t)
  g.Describe("Check Squish bufHash32()", func() {
    g.It("calculates correct Sqush bufhash-32 of an empty string", func() {
      g.Assert(bufHash32("")).Equal(uint32(0x0))
    })
    g.It("calculates correct Sqush bufhash-32 of the string 'Alexander N. Skovpen'", func() {
      g.Assert(bufHash32("Alexander N. Skovpen")).Equal(uint32(0x00efd7be))
    })
  })
}