package areasconfig

import (
	"testing"

	"github.com/askovpen/gossiped/pkg/msgapi"
	. "github.com/franela/goblin"
)

func TestSquishConfig(t *testing.T) {
	msgapi.Areas = msgapi.Areas[:0]
	g := Goblin(t)
	g.Describe("Check SquishConfig", func() {
		g.It("check squishConfigRead()", func() {
			squishConfigRead("../../testdata/squish.cfg")
			g.Assert(len(msgapi.Areas)).Equal(56)
		})
	})
}
