package config

import (
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
  "runtime"
)

type config_s struct {
  Username string
  FidoConfig string
  Log string
}
var (
  Config config_s
  Version string
  PID string
  LongPID string
)

func Read() error {
  Version="0.0.1"
  PID="ATED+"+runtime.GOOS[0:3]+" "+Version
  LongPID="goAtEd "+runtime.GOOS+"/"+runtime.GOARCH+"-"+Version
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
