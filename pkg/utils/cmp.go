package utils

import (
	"strings"
)

// NamesEqual compare names
func NamesEqual(f string, s string) bool {
	return strings.Replace(strings.Trim(f, " "), ".", "", -1) == strings.Replace(strings.Trim(s, " "), ".", "", -1)
}
