package config

import (
	"log"
	"regexp"
	"strings"
)

// ParseColorscheme parses the text definition for a colorscheme and returns the corresponding object
// Colorschemes are made up of color-link statements linking a color group to a list of colors
// For example, color-link keyword (blue,red) makes all keywords have a blue foreground and
// red background
// Todo: Implement to read Golded schemes in future
func ParseColorscheme(text string) ColorScheme {
	parser := regexp.MustCompile(`color-link\s+(\S*)\s+"(.*)"`)

	lines := strings.Split(text, "\n")

	c := make(ColorScheme)

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
				style = StyleDefault
			}
		} else {
			log.Println("Color-link statement is not valid:", line)
		}
	}
	return c
}
