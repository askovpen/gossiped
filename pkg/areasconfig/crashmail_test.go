package areasconfig

import (
	"github.com/askovpen/gossiped/pkg/msgapi"
	. "github.com/franela/goblin"
	"testing"
)

func TestCrashmailConfig(t *testing.T) {
	msgapi.Areas = msgapi.Areas[:0]
	g := Goblin(t)
	g.Describe("Check CrashmailConfig", func() {
		g.It("check crashmailConfigRead()", func() {
			crashmailConfigRead("../../testdata/crashmail.prefs")
			g.Assert(len(msgapi.Areas)).Equal(5)
		})
	})
}
