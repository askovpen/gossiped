package nodelist

import (
        "github.com/askovpen/gossiped/pkg/types"
        "bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type (
        nodelineS       struct {
            Address types.FidoAddr
            BBS string
            City string
            Sysop string
        }
)

// Nodelist contains NodeList
var Nodelist []nodelineS

//Read reads NodeList from the file
func Read(fn string) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
        defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(b)))
        re := regexp.MustCompile(",")
        var z, n, f string
	for scanner.Scan() {
                if scanner.Text()[0] == ';' {
                    continue
                }
		res := re.Split(scanner.Text(), -1)
		if len(res) < 5 {
			continue
		}
		switch strings.ToLower(res[0]) {
		case "zone":
                        z = res[1]
                        n = "0"
                        f = "0"
		case "region":
                        n = res[1]
                        f = "0"
		case "host":
                        n = res[1]
                        f = "0"
		default:
                        f = res[1]
		}
                address := types.AddrFromString(z+":"+n+"/"+f)
                node := nodelineS {
                        Address: *address,
                        BBS: res[2],
                        City: res[3],
                        Sysop: res[4],
                }
                Nodelist = append(Nodelist, node)
	}
	return nil
}
