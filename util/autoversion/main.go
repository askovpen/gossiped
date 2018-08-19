package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func main() {
	cwd, err := filepath.Abs("../..")
	if err != nil {
		panic(err)
	}
	r, err := git.PlainOpen(cwd)
	if err != nil {
		panic(err)
	}
	//head, err := r.Head()
	head, err := r.Head()
	tags, err := r.Tags()
	var tag *plumbing.Reference
	err = tags.ForEach(func(t *plumbing.Reference) error {
		tag = t
		return nil
	})
	//	tag, err := r.TagObject(head.Hash())
	if err != nil {
		panic(err)
	}
	vgo := fmt.Sprintf("package config\n\nvar (\n\tVersion = \"%s-%s\"\n)\n", strings.Split(string(tag.Name()), "-")[1], head.Hash().String()[0:8])
	err = ioutil.WriteFile("../../lib/config/version.go", []byte(vgo), 0644)
	if err != nil {
		panic(err)
	}
}
