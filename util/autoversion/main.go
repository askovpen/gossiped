package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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
	if err != nil {
		panic(err)
	}
	cIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		panic(err)
	}
	cc := 0
	err = cIter.ForEach(func(c *object.Commit) error {
		cc += 1
		return nil
	})
	if err != nil {
		panic(err)
	}
	vgo := fmt.Sprintf("package config\n\nvar (\n\tVersion = \"%s-%d-%s\"\n)\n", strings.Split(string(tag.Name()), "-")[1], cc, head.Hash().String()[0:8])
	err = ioutil.WriteFile("../../lib/config/version.go", []byte(vgo), 0644)
	if err != nil {
		panic(err)
	}
}
