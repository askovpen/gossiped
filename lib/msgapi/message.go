package msgapi

import (
	"errors"
	"fmt"
	"github.com/askovpen/goated/lib/config"
	"github.com/askovpen/goated/lib/types"
	"github.com/askovpen/goated/lib/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Area        string
	AreaID      int
	MsgNum      uint32
	MaxNum      uint32
	DateWritten time.Time
	DateArrived time.Time
	Attrs       []string
	Body        string
	FromAddr    *types.FidoAddr
	ToAddr      *types.FidoAddr
	From        string
	To          string
	Subject     string
	Kludges     map[string]string
}

func (m *Message) ParseRaw() error {
	m.Kludges = make(map[string]string)
	for _, l := range strings.Split(m.Body, "\x0d") {
		if len(l) > 5 && l[0:6] == "\x01INTL " {
			m.Kludges["INTL"] = l[6:]
		} else if len(l) > 5 && l[0:6] == "\x01TOPT " {
			m.Kludges["TOPT"] = l[6:]
		} else if len(l) > 5 && l[0:6] == "\x01FMPT " {
			m.Kludges["FMPT"] = l[6:]
		} else if len(l) > 7 && l[0:8] == "\x01MSGID: " {
			m.Kludges["MSGID:"] = l[8:]
		} else if len(l) > 10 && l[0:11] == "\x20*\x20Origin: " {
			re := regexp.MustCompile("\\d+:\\d+/\\d+\\.*\\d*")
			if len(re.FindStringSubmatch(l)) > 0 {
				m.Kludges["ORIGIN"] = re.FindStringSubmatch(l)[0]
			}
		} else if len(l) > 6 && l[0:7] == "\x01CHRS: " {
			m.Kludges["CHRS"] = strings.ToUpper(strings.Split(l, " ")[1])
		}
	}
	log.Printf("ParseRaw(): %#v", m.Kludges)
	if m.FromAddr == nil {
		if _, ok := m.Kludges["INTL"]; ok {
			m.ToAddr = types.AddrFromString(strings.Split(m.Kludges["INTL"], " ")[0])
			m.FromAddr = types.AddrFromString(strings.Split(m.Kludges["INTL"], " ")[1])
		} else if _, ok := m.Kludges["ORIGIN"]; ok {
			m.FromAddr = types.AddrFromString(m.Kludges["ORIGIN"])
		}
	}
	//log.Printf("%#v", m)
	if m.FromAddr == nil {
		return errors.New("FromAddr not defined")
	}
	if m.ToAddr == nil {
		m.ToAddr = &types.FidoAddr{}
	}
	m.Decode()
	return nil
}

func (m *Message) Encode() {
	enc := strings.Split(config.Config.Chrs, " ")[0]
	m.Body = utils.EncodeCharmap(m.Body, enc)
	m.From = utils.EncodeCharmap(m.From, enc)
	m.To = utils.EncodeCharmap(m.To, enc)
	m.Subject = utils.EncodeCharmap(m.Subject, enc)
}

func (m *Message) Decode() {
	enc := strings.Split(config.Config.Chrs, " ")[0]
	if _, ok := m.Kludges["CHRS"]; ok {
		enc = m.Kludges["CHRS"]
	}
	log.Printf("Decode(): %#v", m.Kludges)
	m.Body = utils.DecodeCharmap(m.Body, enc)
	m.From = utils.DecodeCharmap(m.From, enc)
	m.To = utils.DecodeCharmap(m.To, enc)
	m.Subject = utils.DecodeCharmap(m.Subject, enc)
}

func (m *Message) ToView(showKludges bool) string {
	var nm []string
	re := regexp.MustCompile(">+")
	for _, l := range strings.Split(m.Body, "\x0d") {
		if len(l) > 1 && l[0] == 1 {
			if showKludges {
				nm = append(nm, "\033[30;1m@"+l[1:]+"\033[0m")
			}
		} else if len(l) > 10 && l[0:11] == " * Origin: " {
			nm = append(nm, "\033[37;1m"+l+"\033[0m")
		} else if len(l) > 3 && l[0:4] == "--- " {
			nm = append(nm, "\033[37;1m"+l+"\033[0m")
		} else if len(l) > 3 && l[0:4] == "... " {
			nm = append(nm, "\033[37;1m"+l+"\033[0m")
		} else if len(l) > 8 && l[0:9] == "SEEN-BY: " {
			if showKludges {
				nm = append(nm, "\033[30;1m"+l+"\033[0m")
			}
		} else if ind := re.FindStringIndex(l); ind != nil {
			ind2 := strings.Index(l, "<")
			if (ind2 == -1 || ind2 > ind[1]) && ind[0] < 6 {
				if (ind[1]-ind[0])%2 == 0 {
					nm = append(nm, "\033[37;1m"+l+"\033[0m")
				} else {
					nm = append(nm, "\033[33;1m"+l+"\033[0m")
				}
			} else {
				nm = append(nm, l)
			}
		} else {
			nm = append(nm, l)
		}
	}
	return strings.Join(nm, "\n")
}
func (m *Message) ToEditNewView() (string, int) {
	var nm []string
	p := 0
	r := strings.NewReplacer(
		"@pseudo", m.To,
		"@CFName", strings.Split(m.From, " ")[0])
	for _, l := range config.Template {
		if len(l) > 0 {
			if l[0] == '@' {
				if len(l) > 3 && l[0:4] == "@New" {
					if len(l) == 4 {
						nm = append(nm, "")
					} else {
						nm = append(nm, r.Replace(l[4:]))
					}
				} else if len(l) > 8 && l[0:9] == "@Position" {
					p = len(nm)
					if len(l) == 9 {
						nm = append(nm, "")
					} else {
						nm = append(nm, r.Replace(l[9:]))
					}
				} else if len(l) > 6 && l[0:7] == "@CFName" {
					nm = append(nm, r.Replace(l))
				}
			} else {
				nm = append(nm, r.Replace(l))
			}
		} else {
			nm = append(nm, l)
		}
	}
	nm = append(nm, "\033[37;1m--- "+config.LongPID+"\033[0m")
	nm = append(nm, "\033[37;1m * Origin: "+config.Config.Origin+" ("+m.FromAddr.String()+")\033[0m")
	return strings.Join(nm, "\n"), p

}
func (m *Message) GetQuote() []string {
	var nm []string
	re := regexp.MustCompile(">+")
	from := ""
	for _, l := range strings.Split(m.From, " ") {
		from += string(l[0])
	}
	for _, l := range strings.Split(m.Body, "\x0d") {
		if len(l) > 1 && l[0] == 1 {
			continue
		} else if len(l) > 8 && l[0:9] == "SEEN-BY: " {
			continue
		} else if ind := re.FindStringIndex(l); ind != nil {
			ind2 := strings.Index(l, "<")
			if (ind2 == -1 || ind2 > ind[1]) && ind[0] < 6 {
				if (ind[1]-ind[0])%2 == 0 {
					nm = append(nm, "\033[33;1m"+l[0:ind[0]+1]+">"+l[ind[0]+1:]+"\033[0m")
				} else {
					nm = append(nm, "\033[37;1m"+l[0:ind[0]+1]+">"+l[ind[0]+1:]+"\033[0m")
				}
			} else {
				nm = append(nm, "\033[33;1m "+from+"> "+l+"\033[0m")
			}
		} else {
			nm = append(nm, "\033[33;1m "+from+"> "+l+"\033[0m")
		}
	}
	log.Print(from)
	return nm
}
func (m *Message) ToEditAnswerView(om *Message) (string, int) {
	var nm []string
	p := 0
	r := strings.NewReplacer(
		"@pseudo", m.To,
		"@CFName", strings.Split(m.From, " ")[0],
		"@ODate", om.DateWritten.Format("02 Jan 06"),
		"@OTime", om.DateWritten.Format("15:04:05"),
		"@OName", om.From,
		"@DName", om.To)
	for _, l := range config.Template {
		if len(l) > 0 {
			if l[0] == '@' {
				if len(l) > 15 && l[0:16] == "@Quoted@Position" {
					p = len(nm)
					nm = append(nm, "")
				} else if len(l) > 6 && l[0:7] == "@Quoted" {
					if len(l) == 7 {
						nm = append(nm, "")
					} else {
						nm = append(nm, r.Replace(l[7:]))
					}
				} else if len(l) > 8 && l[0:9] == "@Position" {
					p = len(nm)
					if len(l) == 9 {
						nm = append(nm, "")
					} else {
						nm = append(nm, r.Replace(l[9:]))
					}
				} else if len(l) > 5 && l[0:6] == "@Quote" {
					nm = append(nm, om.GetQuote()...)
				} else if len(l) > 6 && l[0:7] == "@CFName" {
					nm = append(nm, r.Replace(l))
				}
			} else {
				nm = append(nm, r.Replace(l))
			}
		} else {
			nm = append(nm, l)
		}
	}
	nm = append(nm, "\033[37;1m--- "+config.LongPID+"\033[0m")
	nm = append(nm, "\033[37;1m * Origin: "+config.Config.Origin+" ("+m.FromAddr.String()+")\033[0m")
	return strings.Join(nm, "\n"), p
}
func (m *Message) MakeBody() *Message {
	if Areas[m.AreaID].GetType() == EchoAreaTypeNetmail {
		to := m.ToAddr
		top := to.GetPoint()
		to.SetPoint(0)
		from := m.FromAddr
		fromp := from.GetPoint()
		from.SetPoint(0)
		m.Kludges["INTL"] = to.String() + " " + from.String()
		if top > 0 {
			m.Kludges["TOPT"] = strconv.FormatUint(uint64(top), 10)
		}
		if fromp > 0 {
			m.Kludges["FMPT"] = strconv.FormatUint(uint64(fromp), 10)
		}
	}
	m.Kludges["MSGID:"] = fmt.Sprintf("%s %08x", m.FromAddr.String(), uint32(time.Now().Unix()))
	m.Body = strings.Join(strings.Split(m.Body, "\n"), "\x0d")
	m.DateWritten = time.Now()
	m.DateArrived = m.DateWritten
	return m
}
