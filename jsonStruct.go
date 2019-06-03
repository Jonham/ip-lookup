package main

import (
	"encoding/json"
	"io/ioutil"
)

type jsonStruct struct{}

func newJSONStruct() *jsonStruct {
	return &jsonStruct{}
}

// any
func (jst *jsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
