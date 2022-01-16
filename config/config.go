package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Node struct {
		Ip      string
		Port    int
		Id      string
		Timeout int
		ApiPort int `yaml:"api_port"`
	}
	Cluster struct {
		Nodes []string
	}
}

//读取Yaml配置文件,
//并转换成conf对象
func NewConf(application string) *Config {
	//应该是 绝对地址
	var c Config
	yamlFile, err := ioutil.ReadFile(application)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		fmt.Println(err.Error())
	}

	return &c
}
