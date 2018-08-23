package types

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestFidoAddr(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check FidoAddr", func() {
		g.It("check AddrFromString()", func() {
			g.Assert(AddrFromString("2:5020/9696").Equal(&FidoAddr{2, 5020, 9696, 0})).Equal(true)
		})
		g.It("check AddrFromNum()", func() {
			g.Assert(AddrFromNum(2, 5020, 9696, 0).Equal(&FidoAddr{2, 5020, 9696, 0})).Equal(true)
		})
		g.It("check create and compare", func() {
			g.Assert(AddrFromString("2:5020/9696").Equal(AddrFromNum(2, 5020, 9696, 0))).Equal(true)
		})
		g.It("check GetZone()", func() {
			g.Assert((&FidoAddr{2, 5020, 9696, 0}).GetZone()).Equal(uint16(2))
		})
		g.It("check GetNet()", func() {
			g.Assert((&FidoAddr{2, 5020, 9696, 0}).GetNet()).Equal(uint16(5020))
		})
		g.It("check GetNode()", func() {
			g.Assert((&FidoAddr{2, 5020, 9696, 0}).GetNode()).Equal(uint16(9696))
		})
		g.It("check GetPoint()", func() {
			g.Assert((&FidoAddr{2, 5020, 9696, 0}).GetPoint()).Equal(uint16(0))
		})
	})
}
