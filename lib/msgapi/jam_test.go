package msgapi

import (
	. "github.com/franela/goblin"
	"strconv"
	"testing"
)

func TestJamCrc32r(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check Jam crc32r()", func() {
		g.It("calculates correct JAM CRC-32 of an empty string", func() {
			g.Assert(crc32r("")).Equal(uint32(0xffffffff))
		})
		g.It("calculates correct JAM CRC-32 of the string 'Alexander N. Skovpen'", func() {
			g.Assert(crc32r("Alexander N. Skovpen")).Equal(uint32(0x30222bd1))
		})
	})
}

func BenchmarkJamCrc32r(b *testing.B) {
  for n := 0; n < b.N; n++ {
    crc32r(strconv.FormatInt(int64(b.N),10))
  }
}
