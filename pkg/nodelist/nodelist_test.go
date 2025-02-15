package nodelist

import (
        "github.com/askovpen/gossiped/pkg/types"
	. "github.com/franela/goblin"
	"testing"
)

func TestNodelist(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check Nodelist", func() {
		g.It("check nodelist.Read()", func() {
			Read("../../testdata/NODELIST.299")
			g.Assert(len(Nodelist)).Equal(1216)
                        g.Assert(Nodelist[0].Address).Equal(*types.AddrFromString("1:0/0"))
                        g.Assert(Nodelist[0].Sysop).Equal("Nick_Andre")
                        g.Assert(Nodelist[0].BBS).Equal("North_America_(298)")
                        g.Assert(Nodelist[0].City).Equal("Toronto")
                        g.Assert(Nodelist[len(Nodelist)-1].Address).Equal(*types.AddrFromString("4:920/69"))
                        g.Assert(Nodelist[len(Nodelist)-1].Sysop).Equal("John_Dovey")
                        g.Assert(Nodelist[len(Nodelist)-1].BBS).Equal("El_Gato_De_Fuego_BBS_II")
                        g.Assert(Nodelist[len(Nodelist)-1].City).Equal("Pedasi_Panama")
		})
	})
}
