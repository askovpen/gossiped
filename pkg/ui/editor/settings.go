package editor

// DefaultLocalSettings returns the default local settings
// Note that filetype is a local only option
func DefaultLocalSettings() map[string]interface{} {
	return map[string]interface{}{
		"autoindent":     true,
		"colorcolumn":    float64(0),
		"eofnewline":     false,
		"fileformat":     "unix",
		"filetype":       "Unknown",
		"indentchar":     " ",
		"keepautoindent": false,
		"rmtrailingws":   false,
		"savecursor":     false,
		"saveundo":       false,
		"scrollbar":      false,
		"scrollmargin":   float64(3),
		"scrollspeed":    float64(2),
		"softwrap":       true,
		"tabmovement":    false,
		"tabsize":        float64(4),
		"tabstospaces":   false,
		"useprimary":     true,
	}
}
