package server

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	IP               string
	HTTP             uint
	NATStartFromPort uint
}

type NetConfig struct {
	Size           int
	NPeerToConnect int
}

func readFile() ([]byte, error) {
	f, err := os.Open("./config/config.yml")
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	buff.ReadFrom(f)
	return buff.Bytes(), nil
}

func serverConfFromFile() (ServerConfig, error) {
	conf := ServerConfig{}
	b, err := readFile()
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func netConfFromFile() (NetConfig, error) {
	conf := NetConfig{}
	b, err := readFile()
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
