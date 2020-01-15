package utils

import (
	"strings"
)

func NamesEqual(f string, s string) bool {
	if strings.Replace(strings.Trim(f, " "), ".", "", -1) == strings.Replace(strings.Trim(s, " "), ".", "", -1) {
		return true
	}
	return false
}
