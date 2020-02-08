package utils

import (
	//"io/ioutil"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	//"golang.org/x/text/transform"
)

var (
	cDecoder = map[string]*encoding.Decoder{
		"CP866":   charmap.CodePage866.NewDecoder(),
		"+7_FIDO": charmap.CodePage866.NewDecoder(),
		"+7":      charmap.CodePage866.NewDecoder(),
		"IBM866":  charmap.CodePage866.NewDecoder(),
		"CP850":   charmap.CodePage850.NewDecoder(),
		"CP852":   charmap.CodePage852.NewDecoder(),
		"CP848":   charmap.CodePage866.NewDecoder(),
		"CP1250":  charmap.Windows1250.NewDecoder(),
		"CP1251":  charmap.Windows1251.NewDecoder(),
		"CP1252":  charmap.Windows1252.NewDecoder(),
		"CP10000": charmap.Macintosh.NewDecoder(),
		"CP437":   charmap.CodePage437.NewDecoder(),
		"IBMPC":   charmap.CodePage437.NewDecoder(),
		"LATIN-1": charmap.ISO8859_1.NewDecoder(),
		"LATIN-2": charmap.ISO8859_2.NewDecoder(),
		"LATIN-5": charmap.ISO8859_5.NewDecoder(),
		"LATIN-9": charmap.ISO8859_9.NewDecoder(),
	}
	cEncoder = map[string]*encoding.Encoder{
		"CP866":   charmap.CodePage866.NewEncoder(),
		"+7_FIDO": charmap.CodePage866.NewEncoder(),
		"+7":      charmap.CodePage866.NewEncoder(),
		"IBM866":  charmap.CodePage866.NewEncoder(),
		"CP850":   charmap.CodePage850.NewEncoder(),
		"CP852":   charmap.CodePage852.NewEncoder(),
		"CP848":   charmap.CodePage866.NewEncoder(),
		"CP1250":  charmap.Windows1250.NewEncoder(),
		"CP1251":  charmap.Windows1251.NewEncoder(),
		"CP1252":  charmap.Windows1252.NewEncoder(),
		"CP10000": charmap.Macintosh.NewEncoder(),
		"CP437":   charmap.CodePage437.NewEncoder(),
		"IBMPC":   charmap.CodePage437.NewEncoder(),
		"LATIN-1": charmap.ISO8859_1.NewEncoder(),
		"LATIN-2": charmap.ISO8859_2.NewEncoder(),
		"LATIN-5": charmap.ISO8859_5.NewEncoder(),
		"LATIN-9": charmap.ISO8859_9.NewEncoder(),
	}
)

// DecodeCharmap decode string from charmap
func DecodeCharmap(s string, c string) string {
	var dec *encoding.Decoder
	switch chrs := strings.ToUpper(c); chrs {
	case "CP866", "+7_FIDO", "+7", "IBM866", "CP850", "CP852", "CP848", "CP1250", "CP1251", "CP1252", "CP10000", "CP437", "IBMPC", "LATIN-2", "LATIN-5", "LATIN-9":
		dec = cDecoder[chrs]
	case "UTF-8":
		return s
	default:
		dec = cDecoder["LATIN-1"]
	}
	b, err := dec.String(s)
	if err != nil {
		return s
	}
	return b
}

// EncodeCharmap encode string to charmap
func EncodeCharmap(s string, c string) string {
	var enc *encoding.Encoder
	switch c {
	case "CP866", "+7_FIDO", "+7", "IBM866", "CP850", "CP852", "CP848", "CP1250", "CP1251", "CP1252", "CP10000", "CP437", "IBMPC", "LATIN-2", "LATIN-5", "LATIN-9":
		enc = cEncoder[c]
	case "UTF-8":
		return s
	default:
		enc = cEncoder["LATIN-1"]
	}
	out, err := encoding.ReplaceUnsupported(enc).String(s)
	if err != nil {
		panic(err)
	}
	return out
}
