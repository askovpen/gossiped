package config

import (
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
)

type config_s struct {
  Username string
  FidoConfig string
  Log string
}
var (
  Config config_s
)

func Read() error {
  yamlFile, err := ioutil.ReadFile(os.Args[1])
  if err != nil {
    return err
  }
  err = yaml.Unmarshal(yamlFile, &Config)
  if err != nil {
    return err
  }
  return nil
}
