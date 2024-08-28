package editor

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	. "github.com/franela/goblin"
	"github.com/gdamore/tcell/v2"
	"regexp"
	"strings"
	"testing"
)

var (
	lineSplitter = regexp.MustCompile(`[\r\n]+`)
)

func TestProduceColorMapFromConfig(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check ProduceColorMapFromConfig()", func() {
		config.Config.Colors = map[string]config.ColorMap{
			"editor": {
				"comment":  "bold yellow,red",
				"icomment": "bold red",
				"origin":   "bold red",
				"tearline": "bold red",
				"tagline":  "bold red",
				"kludge":   "bold red",
			},
		}
		g.It("override for existing color definitions", func() {
			colors := config.ColorMap{
				"comment":  "bold yellow,red",
				"icomment": "bold white",
				"origin":   "bold white",
				"tearline": "bold white",
				"tagline":  "bold white",
				"kludge":   "bold gray",
			}
			var produced, err = ProduceColorMapFromConfig(ConfigColorArea, &colors)
			g.Assert(err).IsNil()
			g.Assert(produced["kludge"]).Equal("bold red")
		})
		g.It("override for non-existing color definitions", func() {
			colors := config.ColorMap{
				"comment": "bold yellow",
			}
			var errorsGot, errorsExpected = map[string]bool{}, map[string]bool{}
			var area = ConfigColorArea
			var produced, err = ProduceColorMapFromConfig(area, &colors)
			g.Assert(err).IsNotNil()
			for _, i := range lineSplitter.Split(err.Error(), -1) {
				errorsGot[strings.TrimSpace(i)] = true
			}
			for k, _ := range config.Config.Colors[area] {
				if colors[k] == "" {
					errString := "not valid element for area (element: " + k + ", area: " + area + ")"
					errorsExpected[errString] = true
				}
			}
			g.Assert(errorsGot).Equal(errorsExpected)
			g.Assert(produced["kludge"]).Equal("")
		})
		g.It("fallback mode - empty config area", func() {
			colors := config.ColorMap{
				"comment": "bold yellow",
			}
			var produced, err = ProduceColorMapFromConfig("random-area", &colors)
			g.Assert(err).IsNil()
			g.Assert(produced).IsNotNil()
			g.Assert(produced["kludge"]).Equal("")
			g.Assert(produced["comment"]).Equal("bold yellow")
		})
		g.It("invalid config values - empty value ", func() {
			config.Config.Colors = map[string]config.ColorMap{
				"fictive": {
					"comment": "",
				},
			}
			colors := config.ColorMap{
				"comment": "bold yellow",
			}
			var area = "fictive"
			var errorsGot, errorsExpected = map[string]bool{}, map[string]bool{}
			var produced, err = ProduceColorMapFromConfig(area, &colors)
			g.Assert(produced).IsNotNil()
			g.Assert(err).IsNotNil()
			for _, i := range lineSplitter.Split(err.Error(), -1) {
				errorsGot[strings.TrimSpace(i)] = true
			}
			for k, _ := range config.Config.Colors[area] {
				if colors[k] != "" {
					errString := "empty color value (element: " + k + ", area: " + area + ")"
					errorsExpected[errString] = true
				}
			}
			g.Assert(errorsGot).Equal(errorsExpected)
		})
	})
}

func TestProduceColorSchemeFromConfig(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check ProduceColorSchemeFromConfig", func() {
		g.It("fallback mode (config values are absent)", func() {
			config.Config.Colors = map[string]config.ColorMap{
				"editor": nil,
			}
			scheme := ProduceColorSchemeFromConfig()
			for k, _ := range defaultColors {
				g.Assert(scheme[k]).Equal(GetColor(defaultColors[k]))
			}
		})
		g.It("normal mode (config values are present)", func() {
			config.Config.Colors = map[string]config.ColorMap{
				"editor": {
					"comment":  "bold yellow,red",
					"icomment": "bold red",
					"origin":   "bold red",
					"tearline": "bold red",
					"tagline":  "bold red",
					"kludge":   "bold red",
				},
			}
			scheme := ProduceColorSchemeFromConfig()
			for k, _ := range defaultColors {
				g.Assert(scheme[k]).Equal(GetColor(config.Config.Colors["editor"][k]))
			}
		})
	})
}

func TestStringToStyle(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check StringToStyle", func() {
		g.It("check expected success conversion (fg, bg, styles)", func() {
			var testData = map[string]tcell.Style{
				"default":                                tcell.StyleDefault,
				"black,white":                            tcell.Style{}.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite),
				"orange, red":                            tcell.Style{}.Foreground(tcell.ColorOrange).Background(tcell.ColorRed),
				"yellow,  blue":                          tcell.Style{}.Foreground(tcell.ColorYellow).Background(tcell.ColorBlue),
				"underline default,default":              tcell.Style{}.Underline(true),
				"bold default,default":                   tcell.Style{}.Bold(true),
				"bold|reverse|underline default,default": tcell.Style{}.Bold(true).Underline(true).Reverse(true),
				"reverse yellow,red":                     tcell.Style{}.Reverse(true).Foreground(tcell.ColorYellow).Background(tcell.ColorRed),
				"bold 201,114":                           tcell.Style{}.Foreground(tcell.Color201).Background(tcell.Color114).Bold(true),
				"299,294":                                tcell.Style{}.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault),
			}
			for from, to := range testData {
				expected, _ := StringToStyle(from)
				g.Assert(to).Equal(expected, fmt.Sprintf("failed record: %s", from))
			}
		})
		g.It("check unsuccessful conversion", func() {
			var testData = map[string]map[string]bool{
				"careful big, foobar": {
					"unknown foreground color name \"big\"":    true,
					"unknown background color name \"foobar\"": true,
					"unknown style \"careful\"":                true,
				},
			}
			for style, errorsExpected := range testData {
				var errorsGot = map[string]bool{}
				_, err := StringToStyle(style)
				for _, i := range lineSplitter.Split(err.Error(), -1) {
					errorsGot[strings.TrimSpace(i)] = true
				}
				g.Assert(errorsGot).Equal(errorsExpected, fmt.Sprintf("failed record: %s", style))
			}
		})
	})
}
