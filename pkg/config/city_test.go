package config

import (
	"testing"

	"github.com/askovpen/gossiped/pkg/types"
	. "github.com/franela/goblin"
)

func TestGetCity(t *testing.T) {
	err := Read("../../gossiped.example.yml")
	if err != nil {
		panic(err)
	}
	g := Goblin(t)
	g.Describe("Check City", func() {
		g.It("check GetCity() Moscow", func() {
			g.Assert(GetCity(types.AddrFromString("2:5020/9696"))).Equal("Moscow Russia")
		})
		g.It("check GetCity() Gomel", func() {
			g.Assert(GetCity(types.AddrFromString("2:452/28.1"))).Equal("Gomel")
		})
		g.It("check GetCity() Unknown", func() {
			g.Assert(GetCity(types.AddrFromString("2:1020/9696"))).Equal("unknown")
		})
	})
}
