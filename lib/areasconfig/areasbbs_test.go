package areasconfig

import (
	"github.com/askovpen/goated/lib/msgapi"
	. "github.com/franela/goblin"
	"testing"
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
