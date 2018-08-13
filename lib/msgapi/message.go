package msgapi

import (
  "errors"
  "github.com/askovpen/goated/lib/types"
  "github.com/askovpen/goated/lib/utils"
  "log"
  "regexp"
  "strings"
  "time"
)

type Message struct {
  Area  string
  MsgNum  uint32
  MaxNum  uint32
  DateWritten time.Time
  DateArrived time.Time
  Attr uint32
  Body string
  FromAddr *types.FidoAddr
  ToAddr *types.FidoAddr
  From string
  To string
  Subject string
  kludges map[string]string
}

func (m *Message) ParseRaw() error {
  m.kludges=make(map[string]string)
  for _, l := range strings.Split(m.Body, "\x0d") {
    if len(l)>5 && l[0:6]=="\x01INTL " {
      m.kludges["INTL"]=l[6:]
    } else if len(l)>5 && l[0:6]=="\x01TOPT " {
      m.kludges["TOPT"]=l[6:]
    } else if len(l)>5 && l[0:6]=="\x01FMPT " {
      m.kludges["FMPT"]=l[6:]
    }
  }
  log.Printf("%#v", m.kludges)
  if _, ok := m.kludges["INTL"]; ok {
    m.ToAddr=types.AddrFromString(strings.Split(m.kludges["INTL"]," ")[0])
    m.FromAddr=types.AddrFromString(strings.Split(m.kludges["INTL"]," ")[1])
  }
  //log.Printf("%#v", m)
  if m.FromAddr==nil {
    return errors.New("FromAddr not defined")
  }
  if m.ToAddr==nil {
    m.ToAddr=&types.FidoAddr{}
  }
  m.Decode()
  return nil
}

func (m *Message) Decode() {
  enc:="CP866"
  m.Body=utils.DecodeCharmap(m.Body,enc)
  m.From=utils.DecodeCharmap(m.From,enc)
  m.To=utils.DecodeCharmap(m.To,enc)
  m.Subject=utils.DecodeCharmap(m.Subject,enc)
}

func (m *Message) ToView(showKludges bool) string {
  var nm []string
  re := regexp.MustCompile(">+")
  for _, l := range strings.Split(m.Body, "\x0d") {
    
    if len(l)>1 && l[0]==1 {
      if showKludges { nm=append(nm,"\033[30;1m@"+l[1:]+"\033[0m") }
    } else if len(l)>10 && l[0:11]==" * Origin: " {
      nm=append(nm,"\033[37;1m"+l+"\033[0m")
    } else if len(l)>3 && l[0:4]=="--- " {
      nm=append(nm,"\033[37;1m"+l+"\033[0m")
    } else if len(l)>3 && l[0:4]=="... " {
      nm=append(nm,"\033[37;1m"+l+"\033[0m")
    } else if len(l)>8 && l[0:9]=="SEEN-BY: " {
      if showKludges { nm=append(nm,"\033[30;1m"+l+"\033[0m") }
    } else if ind:=re.FindStringIndex(l); ind!=nil {
      ind2:=strings.Index(l,"<")
      if ind2==-1 || ind2>ind[1] {
        if (ind[1]-ind[0])%2==0 {
          nm=append(nm,"\033[37;1m"+l+"\033[0m")
        } else {
          nm=append(nm,"\033[33;1m"+l+"\033[0m")
        }
      } else {
        nm=append(nm,l)
      }
    } else {
      nm=append(nm,l)
    }
  }
  return strings.Join(nm,"\n")
}
