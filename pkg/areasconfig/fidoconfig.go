package areasconfig

import (
	"bufio"
	"errors"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	defaultMsgType msgapi.EchoAreaMsgType
	mp             = map[string]msgapi.EchoAreaType{
		"ECHOAREA":    msgapi.EchoAreaTypeEcho,
		"LOCALAREA":   msgapi.EchoAreaTypeLocal,
		"NETMAILAREA": msgapi.EchoAreaTypeNetmail,
		"DUPEAREA":    msgapi.EchoAreaTypeDupe,
		"BADAREA":     msgapi.EchoAreaTypeBad,
	}
)

func fidoConfigRead(fn string) error {
	defaultMsgType = msgapi.EchoAreaMsgTypeMSG
	return readFile(fn)
}

func checkIncludePath(fn string) (string, error) {
	if _, err := os.Stat(fn); err == nil {
		return fn, nil
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(config.Config.AreaFile.Path), fn)); err == nil {
		return filepath.Join(filepath.Dir(config.Config.AreaFile.Path), fn), nil
	}
	return "", errors.New(fn + " not found")
}

func detectComment(line string) bool {
	str := strings.TrimSpace(line)
	if len(str) > 0 && str[0] == '#' {
		return true
	}
	return false
}

func parseFile(res []string) {
	reEnv := regexp.MustCompile(`\[(.+?)\]`)
	switch tag := strings.ToUpper(res[1]); tag {
	case "INCLUDE":
		readFile(reEnv.ReplaceAllStringFunc(res[2], replaceEnv))
	case "ECHOAREA", "LOCALAREA", "NETMAILAREA", "DUPEAREA", "BADAREA":
		processArea(res[0], mp[tag])
	case "ECHOAREADEFAULTS":
		processDef(res[0])
	}
}

func readFile(fn string) error {
	re := regexp.MustCompile(`(\w+?)\s+(.*)`)
	nfn, err := checkIncludePath(fn)
	if err != nil {
		return err
	}
	file, err := os.Open(nfn)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	for scanner.Scan() {
		if detectComment(scanner.Text()) {
			continue
		}
		res := re.FindStringSubmatch(scanner.Text())
		if len(res) > 2 {
			parseFile(res)
		}
	}
	return nil
}

func replaceEnv(s string) string {
	return os.Getenv(s[1 : len(s)-1])
}

func processDef(areaDef string) {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	res := re.FindAllString(areaDef, -1)
	if len(res) == 2 && strings.EqualFold(res[1], "off") {
		defaultMsgType = msgapi.EchoAreaMsgTypeMSG
		return
	}
	if len(res) < 3 {
		return
	}
	defaultMsgType = getMsgBType(res)
}

func processArea(areaDef string, aType msgapi.EchoAreaType) {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	res := re.FindAllString(areaDef, -1)
	if len(res) < 3 {
		return
	}
	if isPassthrough(res) {
		return
	}
	MsgBType := getMsgBType(res)
	if MsgBType == msgapi.EchoAreaMsgTypeSquish {
		area := &msgapi.Squish{AreaName: res[1], AreaPath: res[2], AreaType: aType}
		msgapi.Areas = append(msgapi.Areas, area)
	} else if MsgBType == msgapi.EchoAreaMsgTypeMSG {
		area := &msgapi.MSG{AreaName: res[1], AreaPath: res[2], AreaType: aType}
		msgapi.Areas = append(msgapi.Areas, area)
	} else if MsgBType == msgapi.EchoAreaMsgTypeJAM {
		area := &msgapi.JAM{AreaName: res[1], AreaPath: res[2], AreaType: aType}
		msgapi.Areas = append(msgapi.Areas, area)
	}
}

func getMsgBType(tokens []string) msgapi.EchoAreaMsgType {
	for i, t := range tokens {
		if strings.EqualFold(t, "-b") {
			if strings.EqualFold(tokens[i+1], "squish") {
				return msgapi.EchoAreaMsgTypeSquish
			} else if strings.EqualFold(tokens[i+1], "jam") {
				return msgapi.EchoAreaMsgTypeJAM
			} else if strings.EqualFold(tokens[i+1], "msg") {
				return msgapi.EchoAreaMsgTypeMSG
			}
			return defaultMsgType
		}
	}
	return defaultMsgType
}

func isPassthrough(tokens []string) bool {
	if strings.EqualFold(tokens[2], "passthrough") {
		return true
	}
	for _, t := range tokens {
		if strings.EqualFold(t, "-pass") {
			return true
		}
	}
	return false
}
