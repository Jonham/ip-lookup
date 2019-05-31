package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type webServiceConf struct {
	Host string
	Port string
}

func main() {
	data, _ := ioutil.ReadFile("config.yml")
	config := webServiceConf{}
	yaml.Unmarshal(data, &config)

	r := setupRouter()

	// 定时提交Get网络下的IP
	startInterval(sendGetNetIP, 0, 5)

	r.Run(config.Host + ":" + config.Port)
}
