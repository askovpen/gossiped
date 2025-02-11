package areasconfig

import (
	"os"
	"testing"

	"github.com/askovpen/gossiped/pkg/msgapi"
	. "github.com/franela/goblin"
)

func TestFidoConfig(t *testing.T) {
	msgapi.Areas = msgapi.Areas[:0]
	g := Goblin(t)
	g.Describe("Check FidoConfig", func() {
		g.It("check replaceEnv()", func() {
			os.Setenv("TESTENV", "PASSED")
			g.Assert(replaceEnv("[TESTENV]")).Equal("PASSED")
		})
		g.It("check fidoConfigRead()", func() {
			fidoConfigRead("../../testdata/hpt.areas")
			g.Assert(len(msgapi.Areas)).Equal(18)
		})
	})
}
