package server

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestLoadConf(t *testing.T) {
	f, err := os.Open("../config/config.yml")
	if err != nil {
		t.Fatal(err)
	}
	var buff bytes.Buffer
	buff.ReadFrom(f)
	b := buff.Bytes()
	//fmt.Println("b:", b)
	var conf Config
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}
