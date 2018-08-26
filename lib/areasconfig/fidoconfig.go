package areasconfig

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
	defaultMsgType msgapi.EchoAreaMsgType
)

func fidoConfigRead(fn string) error {
	defaultMsgType = msgapi.EchoAreaMsgTypeMSG
	readFile(fn)

	if len(msgapi.Areas) == 0 {
		return errors.New("no Areas found")
	}

	sort.Slice(msgapi.Areas, func(i, j int) bool {
		if msgapi.Areas[i].GetType() != msgapi.Areas[j].GetType() {
			return msgapi.Areas[i].GetType() < msgapi.Areas[j].GetType()
		}
		return msgapi.Areas[i].GetName() < msgapi.Areas[j].GetName()
	})
	return nil
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

func readFile(fn string) {
	re := regexp.MustCompile("(\\w+?)\\s+(.*)")
	reEnv := regexp.MustCompile("\\[(.+?)\\]")
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
		res := re.FindStringSubmatch(scanner.Text())
		if len(res) > 2 {
			if strings.EqualFold(res[1], "include") {
				readFile(reEnv.ReplaceAllStringFunc(res[2], replaceEnv))
			} else if strings.EqualFold(res[1], "echoarea") {
				processArea(res[0], msgapi.EchoAreaTypeEcho)
			} else if strings.EqualFold(res[1], "localarea") {
				processArea(res[0], msgapi.EchoAreaTypeLocal)
			} else if strings.EqualFold(res[1], "netmailarea") {
				processArea(res[0], msgapi.EchoAreaTypeNetmail)
			} else if strings.EqualFold(res[1], "EchoAreaDefaults") {
				processDef(res[0])
			}
		}
	}
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
