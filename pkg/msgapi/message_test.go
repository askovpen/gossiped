package msgapi

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestWrapQuoteLine(t *testing.T) {
	g := Goblin(t)
	g.Describe("wrapQuoteLine", func() {
		g.Describe("no wrapping (single segment)", func() {
			g.It("margin == 0 disables wrapping", func() {
				g.Assert(wrapQuoteLine(" YG> ", "some long body that would otherwise wrap", 0)).
					Equal([]string{" YG> some long body that would otherwise wrap"})
			})
			g.It("body fits exactly in avail", func() {
				// prefix " P> " = 4 runes, margin 22 → avail 18, body 18 runes
				g.Assert(wrapQuoteLine(" P> ", "123456789012345678", 22)).
					Equal([]string{" P> 123456789012345678"})
			})
			g.It("body shorter than avail", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16, body 2 runes
				g.Assert(wrapQuoteLine(" P> ", "hi", 20)).
					Equal([]string{" P> hi"})
			})
		})

		g.Describe("word wrap", func() {
			g.It("cuts at last space within window", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "hello world foo bar" (19 runes)
				// window[0..15], runes[16]='b' → roll back to last space in [1..15] at index 15 → "hello world foo"
				g.Assert(wrapQuoteLine(" P> ", "hello world foo bar", 20)).
					Equal([]string{" P> hello world foo", " P> bar"})
			})
			g.It("keeps full window when rune after window is a space", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "hello world test something" (26 runes)
				// window[0..15], runes[16]=' ' → cut at avail → "hello world test"
				// (does not roll back to the earlier space at index 11)
				g.Assert(wrapQuoteLine(" P> ", "hello world test something", 20)).
					Equal([]string{" P> hello world test", " P> something"})
			})
			g.It("trims leading spaces from remainder", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "abcdefghijklmnop    world" (25 runes): window[0..15] has no space → hard break at 16,
				// remainder "    world" trimmed to "world"
				g.Assert(wrapQuoteLine(" P> ", "abcdefghijklmnop    world", 20)).
					Equal([]string{" P> abcdefghijklmnop", " P> world"})
			})
			g.It("produces three segments", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "the quick brown fox jumps over the lazy dog" (43 runes)
				// iter1: window[0..15], last space at 15 → "the quick brown", remainder "fox jumps over the lazy dog"
				// iter2: window[0..15], last space at 14 → "fox jumps over", remainder "the lazy dog"
				// tail: "the lazy dog"
				g.Assert(wrapQuoteLine(" P> ", "the quick brown fox jumps over the lazy dog", 20)).
					Equal([]string{" P> the quick brown", " P> fox jumps over", " P> the lazy dog"})
			})
			g.It("prefix is applied to every segment unchanged", func() {
				// prefix " YG> " = 5 runes, margin 20 → avail 15
				// body "one two three four five six" (27 runes) wraps at least once
				result := wrapQuoteLine(" YG> ", "one two three four five six", 20)
				for _, s := range result {
					// every segment must start with the prefix
					g.Assert(s[0:len(" YG> ")]).Equal(" YG> ")
				}
			})
		})

		g.Describe("hard break", func() {
			g.It("breaks a long word without spaces", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "abcdefghijklmnopqrstuvwxyz" (26 runes, no space) → cut at 16
				g.Assert(wrapQuoteLine(" P> ", "abcdefghijklmnopqrstuvwxyz", 20)).
					Equal([]string{" P> abcdefghijklmnop", " P> qrstuvwxyz"})
			})
			g.It("breaks a long word into multiple chunks", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "abcdefghijklmnopqrstuvwxyzabcdefghijkl" (38 runes, no space)
				// → "abcdefghijklmnop" (16), "qrstuvwxyzabcdef" (16), "ghijkl" (6)
				g.Assert(wrapQuoteLine(" P> ", "abcdefghijklmnopqrstuvwxyzabcdefghijkl", 20)).
					Equal([]string{" P> abcdefghijklmnop", " P> qrstuvwxyzabcdef", " P> ghijkl"})
			})
		})

		g.Describe("UTF-8 / runes", func() {
			g.It("counts runes not bytes for cyrillic body", func() {
				// prefix " ВП> " = 5 runes, margin 20 → avail 15
				// body "здравствуй мир, это проверка" (28 runes)
				// window[0..14], runes[15] == ' ' → cut at avail → "здравствуй мир,"
				g.Assert(wrapQuoteLine(" ВП> ", "здравствуй мир, это проверка", 20)).
					Equal([]string{" ВП> здравствуй мир,", " ВП> это проверка"})
			})
			g.It("counts runes in cyrillic prefix for avail", func() {
				// prefix " ВП> " = 5 runes (7 bytes), margin 20 → avail 15
				// body "проверкапроверкапроверка" (24 runes, no space) → hard break at 15
				g.Assert(wrapQuoteLine(" ВП> ", "проверкапроверкапроверка", 20)).
					Equal([]string{" ВП> проверкапроверк", " ВП> апроверка"})
			})
			g.It("fits cyrillic body without wrapping when on rune boundary", func() {
				// prefix " P> " = 4 runes, margin 20 → avail 16
				// body "привет, мир, при" = 16 runes ≤ avail → single
				g.Assert(wrapQuoteLine(" P> ", "привет, мир, при", 20)).
					Equal([]string{" P> привет, мир, при"})
			})
		})
	})
}
