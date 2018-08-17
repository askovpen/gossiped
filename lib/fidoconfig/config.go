package fidoconfig

import (
	"bufio"
	"errors"
	"github.com/askovpen/goated/lib/config"
	"github.com/askovpen/goated/lib/msgapi"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	defaultMsgType msgapi.EchoAreaType
)

func Read() error {
	defaultMsgType = msgapi.EchoAreaTypeMSG
	readFile(config.Config.FidoConfig)

	if len(msgapi.Areas) == 0 {
		return errors.New("no Areas found")
	}

	sort.Slice(msgapi.Areas, func(i, j int) bool {
		if msgapi.Areas[i].GetType() != msgapi.Areas[j].GetType() {
			if msgapi.Areas[i].GetType() == msgapi.EchoAreaTypeMSG {
				return true
			}
			if msgapi.Areas[j].GetType() == msgapi.EchoAreaTypeMSG {
				return false
			}
		}
		return msgapi.Areas[i].GetName() < msgapi.Areas[j].GetName()
	})
	return nil
}

func checkIncludePath(fn string) (string, error) {
	if _, err := os.Stat(fn); err == nil {
		return fn, nil
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(config.Config.FidoConfig), fn)); err == nil {
		return filepath.Join(filepath.Dir(config.Config.FidoConfig), fn), nil
	}
	return "", errors.New(fn + " not found")
}

func readFile(fn string) {
	re := regexp.MustCompile("(\\w+?)\\s+(.*)")
	nfn, err := checkIncludePath(fn)
	if err != nil {
		log.Print(err)
		return
	}
	file, err := os.Open(nfn)
	if err != nil {
		log.Print(err)
		return
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Print(err)
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(string(b[:])))
	for scanner.Scan() {
		//log.Print(scanner.Text())
		res := re.FindStringSubmatch(scanner.Text())
		if len(res) > 2 {
			//log.Printf("%q",res)
			if strings.EqualFold(res[1], "include") {
				readFile(res[2])
			} else if strings.EqualFold(res[1], "echoarea") {
				processArea(res[0])
			} else if strings.EqualFold(res[1], "localarea") {
				processArea(res[0])
			} else if strings.EqualFold(res[1], "netmailarea") {
				processArea(res[0])
			} else if strings.EqualFold(res[1], "EchoAreaDefaults") {
				processDef(res[0])
			}
		}
	}
}

func processDef(areaDef string) {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	res := re.FindAllString(areaDef, -1)
	if len(res) == 2 && strings.EqualFold(res[1], "off") {
		defaultMsgType = msgapi.EchoAreaTypeMSG
		return
	}
	if len(res) < 3 {
		return
	}
	defaultMsgType = getMsgBType(res)
}

func processArea(areaDef string) {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	res := re.FindAllString(areaDef, -1)
	if len(res) < 3 {
		return
	}
	if isPassthrough(res) {
		return
	}
	MsgBType := getMsgBType(res)
	if MsgBType == msgapi.EchoAreaTypeSquish {
		area := &msgapi.Squish{AreaName: res[1], AreaPath: res[2]}
		msgapi.Areas = append(msgapi.Areas, area)
	} else if MsgBType == msgapi.EchoAreaTypeMSG {
		area := &msgapi.MSG{AreaName: res[1], AreaPath: res[2]}
		msgapi.Areas = append(msgapi.Areas, area)
	} else if MsgBType == msgapi.EchoAreaTypeJAM {
		area := &msgapi.JAM{AreaName: res[1], AreaPath: res[2]}
		msgapi.Areas = append(msgapi.Areas, area)
	}
}

func getMsgBType(tokens []string) msgapi.EchoAreaType {
	for i, t := range tokens {
		if strings.EqualFold(t, "-b") {
			if strings.EqualFold(tokens[i+1], "squish") {
				return msgapi.EchoAreaTypeSquish
			} else if strings.EqualFold(tokens[i+1], "jam") {
				return msgapi.EchoAreaTypeJAM
			} else if strings.EqualFold(tokens[i+1], "msg") {
				return msgapi.EchoAreaTypeMSG
			} else {
				return defaultMsgType
			}
		}
	}
	return defaultMsgType
}

func isPassthrough(tokens []string) bool {
	if tokens[2] == "passthrough" {
		return true
	}
	for _, t := range tokens {
		if strings.EqualFold(t, "-pass") {
			return true
		}
	}
	return false
}
