package areasconfig

import (
	. "github.com/franela/goblin"
	"os"
	"testing"
)

func TestFidoConfig(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check FidoConfig", func() {
		g.It("check replaceEnv()", func() {
			os.Setenv("TESTENV", "PASSED")
			g.Assert(replaceEnv("[TESTENV]")).Equal("PASSED")
		})
	})
}
