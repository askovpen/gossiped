package areasconfig

import (
	"bufio"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func crashmailConfigRead(fn string) error {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(b[:])))
	for scanner.Scan() {
		res := re.FindAllString(scanner.Text(), -1)
		if len(res) != 5 {
			continue
		}
		res[1] = strings.Replace(res[1], "\"", "", -1)
		res[4] = strings.Replace(res[4], "\"", "", -1)
		res[4] = strings.Replace(res[4], "\\\\", "\\", -1)
		aType := msgapi.EchoAreaTypeNone
		if strings.EqualFold(res[0], "area") {
			aType = msgapi.EchoAreaTypeEcho
			if strings.EqualFold(res[1], "bad") {
				aType = msgapi.EchoAreaTypeBad
			}
		} else if strings.EqualFold(res[0], "netmail") {
			aType = msgapi.EchoAreaTypeNetmail
		} else if strings.EqualFold(res[0], "localarea") {
			aType = msgapi.EchoAreaTypeLocal
		} else {
			continue
		}
		if strings.EqualFold(res[3], "jam") {
			area := &msgapi.JAM{AreaName: res[1], AreaPath: res[4], AreaType: aType}
			msgapi.Areas = append(msgapi.Areas, area)
		} else if strings.EqualFold(res[3], "msg") {
			area := &msgapi.MSG{AreaName: res[1], AreaPath: res[4], AreaType: aType}
			msgapi.Areas = append(msgapi.Areas, area)
		}
	}
	return nil
}
