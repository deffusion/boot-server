package server

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerConfig `yaml:"Server"`
	NetConfig    `yaml:"Net"`
}

type ServerConfig struct {
	IP               string `yaml:"IP"`
	HTTP             uint   `yaml:"HTTP"`
	NATStartFromPort uint   `yaml:"NATStartFromPort"`
}

type NetConfig struct {
	Size           int `yaml:"Size"`
	NPeerToConnect int `yaml:"NPeerToConnect"`
}

func configFromFile() (*ServerConfig, *NetConfig, error) {
	f, err := os.Open("./config/config.yml")
	if err != nil {
		return nil, nil, err
	}
	var buff bytes.Buffer
	buff.ReadFrom(f)
	b := buff.Bytes()
	var conf Config
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return nil, nil, err
	}
	return &conf.ServerConfig, &conf.NetConfig, nil
}
