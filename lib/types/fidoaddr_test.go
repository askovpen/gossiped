package types

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestFidoAddr(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check FidoAddr", func() {
		g.It("check AddrFromString()", func() {
			g.Assert(AddrFromString("2:5020/9696.5").Equal(&FidoAddr{2, 5020, 9696, 5})).Equal(true)
			g.Assert(AddrFromString("2:5020")).Equal(AddrFromString("abc"))
		})
		g.It("check AddrFromNum()", func() {
			g.Assert(AddrFromNum(2, 5020, 9696, 0).Equal(&FidoAddr{2, 5020, 9696, 0})).Equal(true)
		})
		g.It("check create and compare", func() {
			g.Assert(AddrFromString("2:5020/9696").Equal(AddrFromNum(2, 5020, 9696, 0))).Equal(true)
			g.Assert(AddrFromString("2:5020/9696.5").Equal(AddrFromNum(2, 5020, 9696, 0))).Equal(false)
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
		g.It("check SetPoint() String()", func() {
			g.Assert((&FidoAddr{2, 5020, 9696, 0}).SetPoint(5).String()).Equal("2:5020/9696.5")
			g.Assert((&FidoAddr{0, 0, 0, 0}).String()).Equal("")
			g.Assert((&FidoAddr{2, 5020, 9696, 0}).String()).Equal("2:5020/9696")
		})
		g.It("check FQDN()", func() {
			_,err:=(&FidoAddr{2, 5020, 9696, 5}).FQDN()
			g.Assert(err.Error()).Equal("point")
			f,err:=(&FidoAddr{2, 5020, 9696, 0}).FQDN()
			g.Assert(f).Equal("f9696.n5020.z2.binkp.net")
		})
	})
}
