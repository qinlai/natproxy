package main

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

type Proxy struct {
	Port  uint16 `json:"ClientPort"`
	Proxy string `json:"ServerProxy"`
}

type Config struct {
	Addr               string   `json:"ServerAddr"`
	Port               uint16   `json:"Port"`
	GateModel          string   `json:"GateModel"`
	AuthKey            string   `json:"AuthKey"`
	ReusePort          bool     `json:"ReusePort"`
	MaxPacket          int      `json:"MaxPacket"`
	MemPoolType        string   `json:"MemPoolType"`
	MemPoolFactor      int      `json:"MemPoolFactor"`
	MemPoolMinChunk    int      `json:"MemPoolMinChunk"`
	MemPoolMaxChunk    int      `json:"MemPoolMaxChunk"`
	MemPoolPageSize    int      `json:"MemPoolPageSize"`
	ClientMaxConn      int      `json:"ClientMaxConn"`
	ClientBufferSize   int      `json:"ClientBufferSize"`
	ClientSendChanSize int      `json:"ClientSendChanSize"`
	Proxys             []*Proxy `json:"Proxys"`
}

var gConf Config

const (
	GateModelsServer string = "gs"
	GateModelsClient string = "gc"
)

func init() {
	loadConfig("./conf.json", &gConf)
}

func loadConfig(file string, o interface{}) {
	reg, e := regexp.Compile("//(.*)")

	if e != nil {
		panic(e)
	}

	raw, err := ioutil.ReadFile(file)

	//json 没有注释，对json进行//注释处理
	raw = reg.ReplaceAll(raw, []byte{})
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(raw, o)

	if err != nil {
		panic(err)
	}
}
