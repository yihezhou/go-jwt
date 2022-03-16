package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type Conf struct {
	JwtKey string `yaml:"jwtKey"`
}

func (conf *Conf) getConf() *Conf  {
	yamlfile, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlfile, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf
}

