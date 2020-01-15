package utils

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestFileExists(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check FileExists()", func() {
		g.It("Check exists", func() {
			g.Assert(FileExists("../../testdata/areas.bbs")).Equal(true)
		})
		g.It("Check not exists", func() {
			g.Assert(FileExists("../../testdata/areas.bbs")).Equal(true)
		})
	})
}
