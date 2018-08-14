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
  if c=="CP866" {
    tr=transform.NewReader(sr, charmap.CodePage866.NewDecoder())
  } else if c=="CP850" {
    tr=transform.NewReader(sr, charmap.CodePage850.NewDecoder())
  } else if c=="CP437" {
    tr=transform.NewReader(sr, charmap.CodePage437.NewDecoder())
  } else {
    tr=transform.NewReader(sr, charmap.ISO8859_1.NewDecoder())
  }
  b, err := ioutil.ReadAll(tr)
  if err!=nil {
    panic(err)
  }
  return string(b)
}
