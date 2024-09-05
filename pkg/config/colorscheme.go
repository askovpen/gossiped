package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
)

const (
	ColorAreaStatusBar     = "statusbar"
	ColorAreaDialog        = "dialog"
	ColorAreaMessageList   = "messageList"
	ColorAreaMessageHeader = "messageHeader"
	ColorAreaAreaList      = "areaList"
	ColorAreaAreaListModal = "areaListModal"
	ColorAreaEditor        = "editor"
	ColorAreaHelp          = "help"
	ColorAreaDefault       = "default"
)
const (
	ColorElementHeader    = "header"
	ColorElementSelection = "selection"
	ColorElementTitle     = "title"
	ColorElementItem      = "item"
	ColorElementHighlight = "highlight"
	ColorElementBorder    = "border"
	ColorElementText      = "text"
	ColorElementPrompt    = "prompt"
	ColorElementWindow    = "window"
)
const (
	StyleUnderline = "underline"
	StyleBold      = "bold"
	StyleReverse   = "reverse"
)

type (
	ColorScheme      map[string]tcell.Style
	ColorSchemeMap   map[string]*ColorScheme
	DefaultColorsMap map[string]*ColorMap
)

var (
	uiColors        = ColorSchemeMap{}
	uiDefaultColors = DefaultColorsMap{
		ColorAreaDefault: {
			ColorElementText: "silver, black",
		},
		ColorAreaAreaList: {
			ColorElementBorder:    "blue",
			ColorElementHeader:    "bold yellow",
			ColorElementTitle:     "bold yellow",
			ColorElementSelection: "white, navy",
			ColorElementItem:      "silver",
			ColorElementHighlight: "bold silver",
			ColorElementPrompt:    "silver",
		},
		ColorAreaAreaListModal: {
			ColorElementBorder:    "red",
			ColorElementHeader:    "bold yellow",
			ColorElementTitle:     "bold yellow",
			ColorElementSelection: "white, navy",
			ColorElementItem:      "silver",
			ColorElementHighlight: "bold silver",
			ColorElementPrompt:    "silver",
		},
		ColorAreaMessageList: {
			ColorElementSelection: "bold white, navy",
			ColorElementHeader:    "bold yellow",
			ColorElementTitle:     "bold yellow",
			ColorElementItem:      "silver",
			ColorElementBorder:    "red",
			ColorElementHighlight: "bold default",
		},
		ColorAreaEditor: {
			"comment":  "bold yellow",
			"icomment": "bold white",
			"origin":   "bold white",
			"tearline": "bold white",
			"tagline":  "bold white",
			"kludge":   "bold gray",
		},
		ColorAreaHelp: {
			ColorElementBorder: "bold blue",
			ColorElementTitle:  "bold yellow",
			ColorElementText:   "default",
		},
		ColorAreaMessageHeader: {
			ColorElementItem:      "silver",
			ColorElementHighlight: "bold silver",
			ColorElementHeader:    "bold silver",
			ColorElementSelection: "silver, navy",
			ColorElementBorder:    "bold blue",
			ColorElementTitle:     "bold yellow",
			ColorElementWindow:    "default",
		},
		ColorAreaDialog: {
			ColorElementItem:      "bold silver",
			ColorElementSelection: "bold silver, navy",
			ColorElementTitle:     "bold yellow",
			ColorElementBorder:    "bold red",
		},
		ColorAreaStatusBar: {
			ColorElementText: "bold white, navy",
		},
	}

	styleToMask = map[string]tcell.AttrMask{
		"B": tcell.AttrBold,
		"U": tcell.AttrUnderline,
		"I": tcell.AttrItalic,
		"L": tcell.AttrBlink,
		"D": tcell.AttrDim,
		"S": tcell.AttrStrikeThrough,
		"R": tcell.AttrReverse,
	}
)

func ProduceColorMapFromConfig(colorArea string, fallbackColors *ColorMap) (*ColorMap, error) {
	var out = make(ColorMap)
	var validKeys = make(map[string]bool)
	if fallbackColors != nil {
		for k, v := range *fallbackColors {
			validKeys[k] = true
			out[k] = v
		}
	}
	var fallback = out
	if Config.Colors[colorArea] == nil || len(Config.Colors[colorArea]) == 0 {
		return &fallback, nil
	}
	var validation error = nil
	for element, colorValue := range Config.Colors[colorArea] {
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
	return &out, validation
}

// ProduceColorSchemeFromConfig
// colorArea: node name in gossiped.yml
// defaultColors: pointer to default ColorMap values
// returns pointer to ColorScheme object
func ProduceColorSchemeFromConfig(colorArea string, defaultColors *ColorMap) *ColorScheme {
	scheme := ColorScheme{}
	colors, err := ProduceColorMapFromConfig(colorArea, defaultColors)
	if err != nil {
		log.Println("Color parse errors: ", err)
	}
	for colorType, colorValue := range *colors {
		scheme[colorType], _ = StringToStyle(colorValue)
	}
	return &scheme
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

// StringToStyle returns a style from a string
// The strings must be in the format "extra foregroundcolor,backgroundcolor"
// The 'extra' can be bold, reverse, or underline
func StringToStyle(str string) (tcell.Style, error) {
	var errStack error
	str = strings.ToLower(strings.TrimSpace(str))

	if len(str) == 0 {
		errStack = errors.New("empty color value")
		return StyleDefault, errStack
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

	var fgColor, bgColor, _ = StyleDefault.Decompose()

	if fg != "" && fg != "default" {
		if _, ok := tcell.ColorNames[fg]; !ok {
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

	style := StyleDefault.Foreground(fgColor).Background(bgColor)
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

// GetColor takes in a syntax group and returns the colorscheme's style for that group
func (colorscheme ColorScheme) GetColor(color string) tcell.Style {
	st := StyleDefault
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

func MaskToStringStyle(attrMask tcell.AttrMask) string {
	style := ""
	for s, v := range styleToMask {
		if (attrMask & v) != 0 {
			style = style + strings.ToLower(s)
		}
	}
	return style
}

// GetColors for config section
func GetColors(section string) *ColorScheme {
	if uiColors[section] == nil {
		uiColors[section] = ProduceColorSchemeFromConfig(section, uiDefaultColors[section])
	}
	return uiColors[section]
}

func GetElementStyle(section string, element string) tcell.Style {
	colors := GetColors(section)
	value, ok := (*colors)[element]
	if !ok {
		return StyleDefault
	}
	return value
}

func FormatTextWithStyle(text string, style tcell.Style) string {
	fg, bg, attrs := style.Decompose()
	return fmt.Sprintf("[%s:%s:%s]%s", fg.String(), bg.String(), MaskToStringStyle(attrs), tview.Escape(text))
}

// initColorAliases()
// append aliases to tcell
func initColorAliases() {
	tcell.ColorNames["cyan"] = tcell.ColorDarkCyan
	tcell.ColorNames["lcyan"] = tcell.ColorLightCyan
	tcell.ColorNames["dcyan"] = tcell.ColorDarkCyan
}

// readColors()
func readColors() error {
	initColorAliases()
	if Config.Colorscheme != "" {
		yamlColors, err := os.ReadFile(Config.Colorscheme)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot read color scheme file: %s", Config.Colorscheme))
		}
		colorsBackup := Config.Colors
		err = yaml.Unmarshal(yamlColors, &Config.Colors)
		if err != nil {
			log.Println(fmt.Sprintf("errors during read of color scheme file: %s", Config.Colorscheme))
			log.Println(fmt.Sprintf("yaml unmarshal errors: %v", err))
			Config.Colors = colorsBackup
		} else {
			log.Println(fmt.Sprintf("color scheme read successfully from file: %s", Config.Colorscheme))
		}
	}
	StyleDefault = GetElementStyle(ColorAreaDefault, ColorElementText)
	StyleDefault.Attributes(tcell.AttrNone)
	return nil
}
