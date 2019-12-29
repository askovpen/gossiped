package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type MTag struct {
	Name string
	dt   time.Time
}

func P(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	out, err := exec.Command("git", "describe", "--tags").Output()
	P(err)
	ver := strings.Trim(string(out), "\n")
	out, err = exec.Command("git", "rev-list", "HEAD", "--count").Output()
	P(err)
	cc := strings.Trim(string(out), "\n")
	out, err = exec.Command("git", "rev-parse", "HEAD").Output()
	P(err)
	h := strings.Trim(string(out), "\n")
	vgo := fmt.Sprintf("package config\n\nvar (\n\tVersion = \"%s-%s-%s\"\n)\n", ver, cc, h[0:8])
	err = ioutil.WriteFile("../../pkg/config/version.go", []byte(vgo), 0644)
	P(err)
}
