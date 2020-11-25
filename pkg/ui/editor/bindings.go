package editor

import (
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

// Actions
const (
	ActionCursorUp            = "CursorUp"
	ActionCursorDown          = "CursorDown"
	ActionCursorPageUp        = "CursorPageUp"
	ActionCursorPageDown      = "CursorPageDown"
	ActionCursorLeft          = "CursorLeft"
	ActionCursorRight         = "CursorRight"
	ActionCursorStart         = "CursorStart"
	ActionCursorEnd           = "CursorEnd"
	ActionInsertNewline       = "InsertNewline"
	ActionInsertSpace         = "InsertSpace"
	ActionBackspace           = "Backspace"
	ActionDelete              = "Delete"
	ActionInsertTab           = "InsertTab"
	ActionCenter              = "Center"
	ActionDeleteLine          = "DeleteLine"
	ActionDeleteToEnd         = "DeleteToEnd"
	ActionIndentSelection     = "IndentSelection"
	ActionOutdentSelection    = "OutdentSelection"
	ActionOutdentLine         = "OutdentLine"
	ActionStart               = "Start"
	ActionEnd                 = "End"
	ActionPageUp              = "PageUp"
	ActionPageDown            = "PageDown"
	ActionHalfPageUp          = "HalfPageUp"
	ActionHalfPageDown        = "HalfPageDown"
	ActionStartOfLine         = "StartOfLine"
	ActionEndOfLine           = "EndOfLine"
	ActionToggleOverwriteMode = "ToggleOverwriteMode"
	ActionEscape              = "Escape"
	ActionInsertEnter         = "InsertEnter"
	ActionUnbindKey           = "UnbindKey"
)

// keyDesc holds the data for a keypress (keycode + modifiers)
type keyDesc struct {
	keyCode   tcell.Key
	modifiers tcell.ModMask
	r         rune
}

// KeyBindings associates key presses with view actions.
type KeyBindings map[keyDesc][]func(*View) bool

// NewKeyBindings returns a new set of keybindings from the given set of binding descriptions.
func NewKeyBindings(bindings map[string]string) KeyBindings {
	return make(KeyBindings).BindKeys(bindings)
}

// BindKey binds a key to a list of actions. If the key is not found or if any action is not found, this function has
// no effect.
func (bindings KeyBindings) BindKey(key string, actions string) KeyBindings {
	k, ok := findKey(key)
	if !ok {
		return bindings
	}

	actionNames := strings.Split(actions, ",")
	if actionNames[0] == "UnbindKey" {
		delete(bindings, k)
		if len(actionNames) == 1 {
			return bindings
		}
		actionNames = append(actionNames[:0], actionNames[1:]...)
	}
	acts := make([]func(*View) bool, 0, len(actionNames))
	for _, actionName := range actionNames {
		action := findAction(actionName)
		if action == nil {
			return bindings
		}
		acts = append(acts, action)
	}

	if len(acts) > 0 {
		bindings[k] = acts
	}

	return bindings
}

// BindKeys binds a set of keys to actions.
func (bindings KeyBindings) BindKeys(keys map[string]string) KeyBindings {
	for k, v := range keys {
		bindings.BindKey(k, v)
	}
	return bindings
}

var bindingActions = map[string]func(*View) bool{
	ActionCursorUp:            (*View).CursorUp,
	ActionCursorDown:          (*View).CursorDown,
	ActionCursorPageUp:        (*View).CursorPageUp,
	ActionCursorPageDown:      (*View).CursorPageDown,
	ActionCursorLeft:          (*View).CursorLeft,
	ActionCursorRight:         (*View).CursorRight,
	ActionCursorStart:         (*View).CursorStart,
	ActionCursorEnd:           (*View).CursorEnd,
	ActionInsertNewline:       (*View).InsertNewline,
	ActionInsertSpace:         (*View).InsertSpace,
	ActionBackspace:           (*View).Backspace,
	ActionDelete:              (*View).Delete,
	ActionInsertTab:           (*View).InsertTab,
	ActionCenter:              (*View).Center,
	ActionDeleteLine:          (*View).DeleteLine,
	ActionDeleteToEnd:         (*View).DeleteToEnd,
	ActionIndentSelection:     (*View).IndentSelection,
	ActionOutdentSelection:    (*View).OutdentSelection,
	ActionOutdentLine:         (*View).OutdentLine,
	ActionStart:               (*View).Start,
	ActionEnd:                 (*View).End,
	ActionPageUp:              (*View).PageUp,
	ActionPageDown:            (*View).PageDown,
	ActionHalfPageUp:          (*View).HalfPageUp,
	ActionHalfPageDown:        (*View).HalfPageDown,
	ActionStartOfLine:         (*View).StartOfLine,
	ActionEndOfLine:           (*View).EndOfLine,
	ActionToggleOverwriteMode: (*View).ToggleOverwriteMode,
	ActionEscape:              (*View).Escape,
	ActionInsertEnter:         (*View).InsertNewline,
}

var bindingKeys = map[string]tcell.Key{
	"Up":             tcell.KeyUp,
	"Down":           tcell.KeyDown,
	"Right":          tcell.KeyRight,
	"Left":           tcell.KeyLeft,
	"UpLeft":         tcell.KeyUpLeft,
	"UpRight":        tcell.KeyUpRight,
	"DownLeft":       tcell.KeyDownLeft,
	"DownRight":      tcell.KeyDownRight,
	"Center":         tcell.KeyCenter,
	"PageUp":         tcell.KeyPgUp,
	"PageDown":       tcell.KeyPgDn,
	"Home":           tcell.KeyHome,
	"End":            tcell.KeyEnd,
	"Insert":         tcell.KeyInsert,
	"Delete":         tcell.KeyDelete,
	"Help":           tcell.KeyHelp,
	"Exit":           tcell.KeyExit,
	"Clear":          tcell.KeyClear,
	"Cancel":         tcell.KeyCancel,
	"Print":          tcell.KeyPrint,
	"Pause":          tcell.KeyPause,
	"Backtab":        tcell.KeyBacktab,
	"F1":             tcell.KeyF1,
	"F2":             tcell.KeyF2,
	"F3":             tcell.KeyF3,
	"F4":             tcell.KeyF4,
	"F5":             tcell.KeyF5,
	"F6":             tcell.KeyF6,
	"F7":             tcell.KeyF7,
	"F8":             tcell.KeyF8,
	"F9":             tcell.KeyF9,
	"F10":            tcell.KeyF10,
	"F11":            tcell.KeyF11,
	"F12":            tcell.KeyF12,
	"F13":            tcell.KeyF13,
	"F14":            tcell.KeyF14,
	"F15":            tcell.KeyF15,
	"F16":            tcell.KeyF16,
	"F17":            tcell.KeyF17,
	"F18":            tcell.KeyF18,
	"F19":            tcell.KeyF19,
	"F20":            tcell.KeyF20,
	"F21":            tcell.KeyF21,
	"F22":            tcell.KeyF22,
	"F23":            tcell.KeyF23,
	"F24":            tcell.KeyF24,
	"F25":            tcell.KeyF25,
	"F26":            tcell.KeyF26,
	"F27":            tcell.KeyF27,
	"F28":            tcell.KeyF28,
	"F29":            tcell.KeyF29,
	"F30":            tcell.KeyF30,
	"F31":            tcell.KeyF31,
	"F32":            tcell.KeyF32,
	"F33":            tcell.KeyF33,
	"F34":            tcell.KeyF34,
	"F35":            tcell.KeyF35,
	"F36":            tcell.KeyF36,
	"F37":            tcell.KeyF37,
	"F38":            tcell.KeyF38,
	"F39":            tcell.KeyF39,
	"F40":            tcell.KeyF40,
	"F41":            tcell.KeyF41,
	"F42":            tcell.KeyF42,
	"F43":            tcell.KeyF43,
	"F44":            tcell.KeyF44,
	"F45":            tcell.KeyF45,
	"F46":            tcell.KeyF46,
	"F47":            tcell.KeyF47,
	"F48":            tcell.KeyF48,
	"F49":            tcell.KeyF49,
	"F50":            tcell.KeyF50,
	"F51":            tcell.KeyF51,
	"F52":            tcell.KeyF52,
	"F53":            tcell.KeyF53,
	"F54":            tcell.KeyF54,
	"F55":            tcell.KeyF55,
	"F56":            tcell.KeyF56,
	"F57":            tcell.KeyF57,
	"F58":            tcell.KeyF58,
	"F59":            tcell.KeyF59,
	"F60":            tcell.KeyF60,
	"F61":            tcell.KeyF61,
	"F62":            tcell.KeyF62,
	"F63":            tcell.KeyF63,
	"F64":            tcell.KeyF64,
	"CtrlSpace":      tcell.KeyCtrlSpace,
	"CtrlA":          tcell.KeyCtrlA,
	"CtrlB":          tcell.KeyCtrlB,
	"CtrlC":          tcell.KeyCtrlC,
	"CtrlD":          tcell.KeyCtrlD,
	"CtrlE":          tcell.KeyCtrlE,
	"CtrlF":          tcell.KeyCtrlF,
	"CtrlG":          tcell.KeyCtrlG,
	"CtrlH":          tcell.KeyCtrlH,
	"CtrlI":          tcell.KeyCtrlI,
	"CtrlJ":          tcell.KeyCtrlJ,
	"CtrlK":          tcell.KeyCtrlK,
	"CtrlL":          tcell.KeyCtrlL,
	"CtrlM":          tcell.KeyCtrlM,
	"CtrlN":          tcell.KeyCtrlN,
	"CtrlO":          tcell.KeyCtrlO,
	"CtrlP":          tcell.KeyCtrlP,
	"CtrlQ":          tcell.KeyCtrlQ,
	"CtrlR":          tcell.KeyCtrlR,
	"CtrlS":          tcell.KeyCtrlS,
	"CtrlT":          tcell.KeyCtrlT,
	"CtrlU":          tcell.KeyCtrlU,
	"CtrlV":          tcell.KeyCtrlV,
	"CtrlW":          tcell.KeyCtrlW,
	"CtrlX":          tcell.KeyCtrlX,
	"CtrlY":          tcell.KeyCtrlY,
	"CtrlZ":          tcell.KeyCtrlZ,
	"CtrlLeftSq":     tcell.KeyCtrlLeftSq,
	"CtrlBackslash":  tcell.KeyCtrlBackslash,
	"CtrlRightSq":    tcell.KeyCtrlRightSq,
	"CtrlCarat":      tcell.KeyCtrlCarat,
	"CtrlUnderscore": tcell.KeyCtrlUnderscore,
	"Tab":            tcell.KeyTab,
	"Esc":            tcell.KeyEsc,
	"Escape":         tcell.KeyEscape,
	"Enter":          tcell.KeyEnter,
	"Backspace":      tcell.KeyBackspace2,
	"OldBackspace":   tcell.KeyBackspace,

	// I renamed these keys to PageUp and PageDown but I don't want to break someone's keybindings
	"PgUp":   tcell.KeyPgUp,
	"PgDown": tcell.KeyPgDn,
}

// DefaultKeyBindings DefaultKeyBindings
var DefaultKeyBindings KeyBindings

// InitBindings initializes the keybindings for micro
func init() {
	DefaultKeyBindings = NewKeyBindings(map[string]string{
		"Up":        ActionCursorUp,
		"Down":      ActionCursorDown,
		"Right":     ActionCursorRight,
		"Left":      ActionCursorLeft,
		"Enter":     ActionInsertNewline,
		"CtrlH":     ActionBackspace,
		"Backspace": ActionBackspace,
		"Tab":       ActionIndentSelection + "," + ActionInsertTab,
		"Backtab":   ActionOutdentSelection + "," + ActionOutdentLine,
		"CtrlY":     ActionDeleteLine,
		"CtrlK":     ActionDeleteToEnd,
		"Home":      ActionStartOfLine,
		"End":       ActionEndOfLine,
		"CtrlHome":  ActionCursorStart,
		"CtrlEnd":   ActionCursorEnd,
		"PageUp":    ActionCursorPageUp,
		"PageDown":  ActionCursorPageDown,
		"Delete":    ActionDelete,
		"Insert":    ActionToggleOverwriteMode,
		"Esc":       ActionEscape,
		"F2":        ActionEscape,
	})
}

// findKey will find binding Key 'b' using string 'k'
func findKey(k string) (b keyDesc, ok bool) {
	modifiers := tcell.ModNone

	// First, we'll strip off all the modifiers in the name and add them to the
	// ModMask
modSearch:
	for {
		switch {
		case strings.HasPrefix(k, "-"):
			// We optionally support dashes between modifiers
			k = k[1:]
		case strings.HasPrefix(k, "Ctrl") && k != "CtrlH":
			// CtrlH technically does not have a 'Ctrl' modifier because it is really backspace
			k = k[4:]
			modifiers |= tcell.ModCtrl
		case strings.HasPrefix(k, "Alt"):
			k = k[3:]
			modifiers |= tcell.ModAlt
		case strings.HasPrefix(k, "Shift"):
			k = k[5:]
			modifiers |= tcell.ModShift
		default:
			break modSearch
		}
	}

	if len(k) == 0 {
		return keyDesc{}, false
	}

	// Control is handled specially, since some character codes in bindingKeys
	// are different when Control is depressed. We should check for Control keys
	// first.
	if modifiers&tcell.ModCtrl != 0 {
		// see if the key is in bindingKeys with the Ctrl prefix.
		k = string(unicode.ToUpper(rune(k[0]))) + k[1:]
		if code, ok := bindingKeys["Ctrl"+k]; ok {
			// It is, we're done.
			return keyDesc{
				keyCode:   code,
				modifiers: modifiers,
				r:         0,
			}, true
		}
	}

	// See if we can find the key in bindingKeys
	if code, ok := bindingKeys[k]; ok {
		return keyDesc{
			keyCode:   code,
			modifiers: modifiers,
			r:         0,
		}, true
	}

	// If we were given one character, then we've got a rune.
	if len(k) == 1 {
		return keyDesc{
			keyCode:   tcell.KeyRune,
			modifiers: modifiers,
			r:         rune(k[0]),
		}, true
	}

	// We don't know what happened.
	return keyDesc{}, false
}

// findAction will find 'action' using string 'v'
func findAction(v string) (action func(*View) bool) {
	action = bindingActions[v]
	return action
}
