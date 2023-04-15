package utils

import (
	"encoding/json"
	"io/ioutil"
)

type GlobalObj struct {
	Name string

	Ip   string
	Port uint32

	MaxDataLen      uint32
	MaxMsgRWChanLen uint32

	PoolSize      uint32
	TaskQueueSize uint32
}

var Global_obj *GlobalObj

func init() {
	Global_obj = &GlobalObj{
		Name: "zinx",
		Ip:   "127.0.0.1",
		Port: 8999,

		MaxDataLen:      4096,
		MaxMsgRWChanLen: 1,

		PoolSize:      10,
		TaskQueueSize: 20,
	}

	buf, err := ioutil.ReadFile("config/config.json")
	// invalid json file or not exist...
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, Global_obj)
	if err != nil {
		panic(err.Error())
	}
}
