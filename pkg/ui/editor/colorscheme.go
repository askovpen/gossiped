package editor

import (
	"errors"
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/gdamore/tcell/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Colorscheme is a map from string to style -- it represents a colorscheme
type Colorscheme map[string]tcell.Style

const ConfigColorArea = "editor"

const (
	StyleUnderline = "underline"
	StyleBold      = "bold"
	StyleReverse   = "reverse"
)

var (
	// The current default colorscheme
	colorscheme Colorscheme
	// The default cell style
	defStyle tcell.Style
	// Default colors
	defaultColors = config.ColorMap{
		"comment":  "bold yellow",
		"icomment": "bold white",
		"origin":   "bold white",
		"tearline": "bold white",
		"tagline":  "bold white",
		"kludge":   "bold gray",
	}
)

// GetColor takes in a syntax group and returns the colorscheme's style for that group
func GetColor(color string) tcell.Style {
	return colorscheme.GetColor(color)
}

// GetColor takes in a syntax group and returns the colorscheme's style for that group
func (colorscheme Colorscheme) GetColor(color string) tcell.Style {
	st := defStyle
	if color == "" {
		return st
	}
	groups := strings.Split(color, ".")
	if len(groups) > 1 {
		curGroup := ""
		for i, g := range groups {
			if i != 0 {
				curGroup += "."
			}
			curGroup += g
			if style, ok := colorscheme[curGroup]; ok {
				st = style
			}
		}
	} else if style, ok := colorscheme[color]; ok {
		st = style
	} else {
		st, _ = StringToStyle(color)
	}

	return st
}

// init picks and initializes the colorscheme when micro starts
func init() {
	colorscheme = make(Colorscheme)
	defStyle = tcell.StyleDefault.
		Foreground(tcell.ColorDefault).
		Background(tcell.ColorDefault)
}

// SetDefaultColorscheme sets the current default colorscheme for new Views.
func SetDefaultColorscheme(scheme Colorscheme) {
	colorscheme = scheme
}

// ParseColorscheme parses the text definition for a colorscheme and returns the corresponding object
// Colorschemes are made up of color-link statements linking a color group to a list of colors
// For example, color-link keyword (blue,red) makes all keywords have a blue foreground and
// red background
func ParseColorscheme(text string) Colorscheme {
	parser := regexp.MustCompile(`color-link\s+(\S*)\s+"(.*)"`)

	lines := strings.Split(text, "\n")

	c := make(Colorscheme)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" ||
			strings.TrimSpace(line)[0] == '#' {
			// Ignore this line
			continue
		}

		matches := parser.FindSubmatch([]byte(line))
		if len(matches) == 3 {
			link := string(matches[1])
			colors := string(matches[2])

			style, _ := StringToStyle(colors)
			c[link] = style

			if link == "default" {
				defStyle = style
			}
		} else {
			log.Println("Color-link statement is not valid:", line)
		}
	}

	return c
}

// StringToStyle returns a style from a string
// The strings must be in the format "extra foregroundcolor,backgroundcolor"
// The 'extra' can be bold, reverse, or underline
func StringToStyle(str string) (tcell.Style, error) {
	var errStack error
	str = strings.ToLower(strings.TrimSpace(str))

	if len(str) == 0 {
		errStack = errors.New("empty color value")
		return tcell.StyleDefault, errStack
	}

	var fg, bg string
	var split = strings.Split(str, ",")
	if len(split) > 1 {
		fg, bg = split[0], split[1]
	} else {
		fg = split[0]
	}
	fg = strings.TrimSpace(fg)
	bg = strings.TrimSpace(bg)

	var styles = ""
	var splitFg = strings.Split(fg, " ")
	if len(splitFg) > 1 {
		styles = strings.TrimSpace(splitFg[0])
		fg = strings.TrimSpace(splitFg[1])
	} else {
		fg = strings.TrimSpace(splitFg[0])
	}

	var fgColor, bgColor, _ = defStyle.Decompose()

	if fg != "" && fg != "default" {
		if _, ok := tcell.ColorNames[fg]; !ok && fg != "default" {
			errStack = errors.Join(errStack, errors.New(fmt.Sprintf("unknown foreground color name \"%s\"", fg)))
		}
		fgColor = StringToColor(fg)
	}
	if bg != "" && bg != "default" {
		if _, ok := tcell.ColorNames[bg]; !ok {
			errStack = errors.Join(errStack, errors.New(fmt.Sprintf("unknown background color name \"%s\"", bg)))
		}
		bgColor = StringToColor(bg)
	}

	style := defStyle.Foreground(fgColor).Background(bgColor)
	var splitStyles = strings.Split(styles, "|")
	for _, v := range splitStyles {
		v = strings.TrimSpace(v)
		if v == StyleReverse {
			style = style.Reverse(true)
		} else if v == StyleUnderline {
			style = style.Underline(true)
		} else if v == StyleBold {
			style = style.Bold(true)
		} else if v != "" {
			errStack = errors.Join(errStack, errors.New(fmt.Sprintf("unknown style \"%s\"", v)))
		}
	}
	return style, errStack
}

// StringToColor returns a tcell color from a string representation of a color
func StringToColor(str string) tcell.Color {
	if num, err := strconv.Atoi(str); err == nil {
		if num > 255 || num < 0 {
			return tcell.ColorDefault
		}
		return tcell.PaletteColor(num)
	}
	return tcell.GetColor(str)
}

func ProduceColorMapFromConfig(colorArea string, fallbackColors *config.ColorMap) (config.ColorMap, error) {
	var out = make(config.ColorMap)
	var validKeys = make(map[string]bool)
	for k, v := range *fallbackColors {
		validKeys[k] = true
		out[k] = v
	}
	var fallback = out
	if config.Config.Colors[colorArea] == nil || len(config.Config.Colors[colorArea]) == 0 {
		return fallback, nil
	}
	var validation error = nil
	for element, colorValue := range config.Config.Colors[colorArea] {
		colorValue = strings.ToLower(strings.TrimSpace(colorValue))
		if !validKeys[element] {
			validation = errors.Join(
				validation,
				errors.New("not valid element for area (element: "+element+", area: "+colorArea+")"),
			)
			continue
		}

		if _, err := StringToStyle(colorValue); err != nil {
			validation = errors.Join(
				validation,
				errors.New(err.Error()+" (element: "+element+", area: "+colorArea+")"),
			)
			continue
		}
		out[element] = colorValue
	}
	return out, validation
}

func ProduceColorSchemeFromConfig() Colorscheme {
	var scheme = Colorscheme{}
	colors, _ := ProduceColorMapFromConfig(ConfigColorArea, &defaultColors)
	for colorType, colorValue := range colors {
		scheme[colorType], _ = StringToStyle(colorValue)
	}
	return scheme
}
