package utils

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

func DecodeCharmap(s string, c string) string {
	sr := strings.NewReader(s)
	var tr *transform.Reader
	switch c {
	case "CP866", "+7_FIDO", "+7":
		tr = transform.NewReader(sr, charmap.CodePage866.NewDecoder())
	case "CP850":
		tr = transform.NewReader(sr, charmap.CodePage850.NewDecoder())
	case "CP852":
		tr = transform.NewReader(sr, charmap.CodePage852.NewDecoder())
	case "CP848":
		// CP848 IBM codepage 848 (Cyrillic Ukrainian) -> to be added as XUserDefined
		tr = transform.NewReader(sr, charmap.CodePage866.NewDecoder())
	case "CP1250":
		tr = transform.NewReader(sr, charmap.Windows1250.NewDecoder())
	case "CP1251":
		tr = transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	case "CP1252":
		tr = transform.NewReader(sr, charmap.Windows1252.NewDecoder())
	case "CP10000":
		tr = transform.NewReader(sr, charmap.Macintosh.NewDecoder())
	case "CP437", "IBMPC":
		tr = transform.NewReader(sr, charmap.CodePage437.NewDecoder())
	case "LATIN-2":
		tr = transform.NewReader(sr, charmap.ISO8859_2.NewDecoder())
	case "LATIN-5":
		tr = transform.NewReader(sr, charmap.ISO8859_5.NewDecoder())
	case "LATIN-9":
		tr = transform.NewReader(sr, charmap.ISO8859_9.NewDecoder())
	case "UTF-8":
		return string(s)
	default:
		tr = transform.NewReader(sr, charmap.ISO8859_1.NewDecoder())
	}
	b, err := ioutil.ReadAll(tr)
	if err != nil {
		panic(err)
	}
	return string(b)
}
