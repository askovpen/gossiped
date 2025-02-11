package areasconfig

import (
	"testing"

	"github.com/askovpen/gossiped/pkg/msgapi"
	. "github.com/franela/goblin"
)

func TestAreasbbsConfig(t *testing.T) {
	msgapi.Areas = msgapi.Areas[:0]
	g := Goblin(t)
	g.Describe("Check AreasbbsConfig", func() {
		g.It("check areasbbsConfigRead()", func() {
			areasbbsConfigRead("../../testdata/areas.bbs")
			g.Assert(len(msgapi.Areas)).Equal(25)
		})
	})
}
