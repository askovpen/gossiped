package utils

import (
  "golang.org/x/text/encoding/charmap"
  "golang.org/x/text/transform"
  "io/ioutil"
  "strings"
)

func DecodeCharmap(s string, c string) string {
  sr:=strings.NewReader(s)
  var tr *transform.Reader
  switch {
    case c=="CP866" || c=="+7_FIDO" || c=="+7" :
      tr=transform.NewReader(sr, charmap.CodePage866.NewDecoder())
    case c=="CP850" :
      tr=transform.NewReader(sr, charmap.CodePage850.NewDecoder())
    case c=="CP852" :
      tr=transform.NewReader(sr, charmap.CodePage852.NewDecoder())
    case c=="CP848" :
      // CP848 IBM codepage 848 (Cyrillic Ukrainian) -> to be added as XUserDefined
      tr=transform.NewReader(sr, charmap.CodePage866.NewDecoder())
    case c=="CP1250" :
      tr=transform.NewReader(sr, charmap.Windows1250.NewDecoder())
    case c=="CP1251" :
      tr=transform.NewReader(sr, charmap.Windows1251.NewDecoder())
    case c=="CP1252" :
      tr=transform.NewReader(sr, charmap.Windows1252.NewDecoder())
    case c=="CP10000" :
      tr=transform.NewReader(sr, charmap.Macintosh.NewDecoder())
    case c=="CP437" || c=="IBMPC" :
      tr=transform.NewReader(sr, charmap.CodePage437.NewDecoder())
    case c=="LATIN-2" :
      tr=transform.NewReader(sr, charmap.ISO8859_2.NewDecoder())
    case c=="LATIN-5" :
      tr=transform.NewReader(sr, charmap.ISO8859_5.NewDecoder())
    case c=="LATIN-9" :
      tr=transform.NewReader(sr, charmap.ISO8859_9.NewDecoder())
    case c=="UTF-8" :
      return string(s)
    default :
      tr=transform.NewReader(sr, charmap.ISO8859_1.NewDecoder())
  }
  b, err := ioutil.ReadAll(tr)
  if err!=nil {
    panic(err)
  }
  return string(b)
}
