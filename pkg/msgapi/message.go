package msgapi

import (
	// "errors"
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/types"
	"github.com/askovpen/gossiped/pkg/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MessageListItem struct
type MessageListItem struct {
	MsgNum      uint32
	From        string
	To          string
	Subject     string
	DateWritten time.Time
}

// Message struct
type Message struct {
	Area        string
	AreaID      int
	MsgNum      uint32
	MaxNum      uint32
	DateWritten time.Time
	DateArrived time.Time
	Attrs       []string
	ReplyTo     uint32
	Replies     []uint32
	Body        string
	FromAddr    *types.FidoAddr
	ToAddr      *types.FidoAddr
	From        string
	To          string
	Subject     string
	Kludges     map[string]string
	Corrupted   bool
}

// ParseRaw parse raw msg
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
	//log.Printf("ParseRaw(): %#v", m.Kludges)
	if m.FromAddr == nil {
		if _, ok := m.Kludges["INTL"]; ok {
			m.ToAddr = types.AddrFromString(strings.Split(m.Kludges["INTL"], " ")[0])
			m.FromAddr = types.AddrFromString(strings.Split(m.Kludges["INTL"], " ")[1])
		} else if _, ok := m.Kludges["ORIGIN"]; ok {
			m.FromAddr = types.AddrFromString(m.Kludges["ORIGIN"])
		}
	}
	if m.FromAddr == nil {
		//return errors.New("FromAddr not defined")
		m.Corrupted = true
		m.FromAddr = &types.FidoAddr{}
	}
	if m.ToAddr == nil {
		m.ToAddr = &types.FidoAddr{}
	}
	if Areas[m.AreaID].GetType() == EchoAreaTypeNetmail {
		if _, ok := m.Kludges["FMPT"]; ok {
			a, err := strconv.ParseUint(m.Kludges["FMPT"], 10, 16)
			if err == nil {
				m.FromAddr.SetPoint(uint16(a))
			}
		}
		if _, ok := m.Kludges["TOPT"]; ok {
			a, err := strconv.ParseUint(m.Kludges["TOPT"], 10, 16)
			if err == nil {
				m.ToAddr.SetPoint(uint16(a))
			}
		}
	}
	m.Decode()
	return nil
}

func (m *Message) parseTabs(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\x09' {
			ts := 8 - (i % 8)
			for j := 0; j < ts; j++ {
				s = s[:i] + " " + s[i:]
				i++
			}
		}
	}
	return s
}

// Encode charset
func (m *Message) Encode() {
	enc := strings.Split(config.Config.Chrs.Default, " ")[0]
	if Areas[m.AreaID].GetChrs() != "" {
		enc = strings.Split(Areas[m.AreaID].GetChrs(), " ")[0]
	}
	m.Body = utils.EncodeCharmap(m.Body, enc)
	m.From = utils.EncodeCharmap(m.From, enc)
	m.To = utils.EncodeCharmap(m.To, enc)
	m.Subject = utils.EncodeCharmap(m.Subject, enc)
}

// Decode charset
func (m *Message) Decode() {
	enc := strings.Split(config.Config.Chrs.Default, " ")[0]
	if _, ok := m.Kludges["CHRS"]; ok {
		enc = m.Kludges["CHRS"]
		if enc == "IBMPC" {
			enc = config.Config.Chrs.IBMPC
		}
	}
	//log.Printf("Decode(): %#v", m.Kludges)
	m.Body = utils.DecodeCharmap(m.Body, enc)
	m.From = utils.DecodeCharmap(m.From, enc)
	m.To = utils.DecodeCharmap(m.To, enc)
	m.Subject = utils.DecodeCharmap(m.Subject, enc)
}

// ToView export view
func (m *Message) ToView(showKludges bool) string {
	var nm []string
	re := regexp.MustCompile(">+")
	for _, l := range strings.Split(m.Body, "\x0d") {
		l = m.parseTabs(l)
		if len(l) > 1 && l[0] == 1 {
			if showKludges {
				nm = append(nm, "[::b][black]@"+l[1:])
			}
		} else if len(l) > 10 && l[0:11] == " * Origin: " {
			nm = append(nm, "[::b]"+l)
		} else if len(l) > 3 && l[0:4] == "--- " {
			nm = append(nm, "[::b]"+l)
		} else if len(l) > 3 && l[0:4] == "... " {
			nm = append(nm, "[::b]"+l)
		} else if len(l) > 8 && l[0:9] == "SEEN-BY: " {
			if showKludges {
				nm = append(nm, "[::b][black]"+l)
			}
		} else if ind := re.FindStringIndex(l); ind != nil {
			ind2 := strings.Index(l, "<")
			if (ind2 == -1 || ind2 > ind[1]) && ind[0] < 6 {
				if (ind[1]-ind[0])%2 == 0 {
					nm = append(nm, "[::b]"+l)
				} else {
					nm = append(nm, "[::b][yellow]"+l)
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

// ToEditNewView export view
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
	nm = append(nm, "--- "+config.Config.Tearline)
	nm = append(nm, " * Origin: "+config.Config.Origin+" ("+m.FromAddr.String()+")")
	log.Printf("pp: %d", p)
	return strings.Join(nm, "\n"), p

}

// GetForward get forward
func (m *Message) GetForward() []string {
	reO := regexp.MustCompile("^ \\* Origin: ")
	reT := regexp.MustCompile("^--- ")
	re := regexp.MustCompile(">+")
	var nm []string
	for _, l := range strings.Split(m.Body, "\x0d") {
		if len(l) > 0 && l[0] == 1 {
			continue
		} else if len(l) > 8 && l[0:9] == "SEEN-BY: " {
			continue
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
			l = reO.ReplaceAllString(l, " + Origin: ")
			l = reT.ReplaceAllString(l, "-+- ")
			nm = append(nm, l)
		}
	}
	return nm
}

// GetQuote get quote
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
					nm = append(nm, l[0:ind[0]+1]+">"+l[ind[0]+1:])
				} else {
					nm = append(nm, l[0:ind[0]+1]+">"+l[ind[0]+1:])
				}
			} else {
				nm = append(nm, " "+from+"> "+l)
			}
		} else {
			nm = append(nm, " "+from+"> "+l)
		}
	}
	//log.Print(from)
	return nm
}

// ToEditAnswerView export view
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
	nm = append(nm, "--- "+config.Config.Tearline)
	nm = append(nm, " * Origin: "+config.Config.Origin+" ("+m.FromAddr.String()+")")
	return strings.Join(nm, "\n"), p
}

// ToEditForwardView export view
func (m *Message) ToEditForwardView(om *Message) (string, int) {
	var nm []string
	p := 0
	r := strings.NewReplacer(
		"@pseudo", m.To,
		"@CFName", strings.Split(m.From, " ")[0],
		"@ODate", om.DateWritten.Format("02 Jan 06"),
		"@OTime", om.DateWritten.Format("15:04:05"),
		"@OName", om.From,
		"@OAddr", om.FromAddr.String(),
		"@DName", om.To,
		"@OEcho", Areas[om.AreaID].GetName(),
		"@Subject", om.Subject,
		"@CAddr", config.Config.Address.String(),
		"@CName", config.Config.Username)
	for _, l := range config.Template {
		if len(l) > 0 {
			if l[0] == '@' {
				if len(l) > 7 && l[0:8] == "@Forward" {
					if len(l) == 8 {
						nm = append(nm, "")
					} else {
						nm = append(nm, r.Replace(l[8:]))
					}
				} else if len(l) > 8 && l[0:9] == "@Position" {
					p = len(nm)
					if len(l) == 9 {
						nm = append(nm, "")
					} else {
						nm = append(nm, r.Replace(l[9:]))
					}
				} else if len(l) > 7 && l[0:8] == "@Message" {
					nm = append(nm, om.GetForward()...)
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
	nm = append(nm, "\033[37;1m--- "+config.Config.Tearline+"\033[0m")
	nm = append(nm, "\033[37;1m * Origin: "+config.Config.Origin+" ("+m.FromAddr.String()+")\033[0m")
	return strings.Join(nm, "\n"), p
}

// MakeBody make body
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
	//time.Sleep(time.Second)
	return m
}
