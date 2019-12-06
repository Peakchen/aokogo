package testconfig

import (
	"common/Log"
	"common/utls"
	"encoding/json"
	"simulate/AutoTest/gameMsg"
)

type tMapInterface map[string]interface{}

type TMsgDsc struct {
	SubId  int32                  `json:"subId"`
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

type TMsgOper struct {
	Request  *TMsgDsc `json:"request"`
	Response *TMsgDsc `json:"response"`
}

type TMsgDetial struct {
	Data []*TMsgOper
}

type tArrMsgDetial []*TMsgDetial

func (this *TMsgDetial) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		err = nil
		return
	}

	dst := &TMsgOper{}
	err = json.Unmarshal(data, &dst)
	if err != nil {
		Log.Error("unmarshal nsg fail, err: ", err)
		return
	}

	this.Data = append(this.Data, dst)
	return
}

type TMsg struct {
	Msg []*TMsgDetial `json:"msg"`
}

func (this *TMsg) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		err = nil
		return
	}

	dst := &tArrMsgDetial{}
	err = json.Unmarshal(data, &dst)
	if err != nil {
		Log.Error("unmarshal nsg fail, err: ", err)
		return
	}

	this.Msg = *dst
	return
}

type TConfig4Test struct {
	ConnAddr string `json:"connAddr"`
	Module   string `json:"module"`
	Route    int32  `json:"route"`
	MainId   int32  `json:"mainId"`
	Msg      *TMsg  `json:"msg"`
}

type TArrConfig4Test []*TConfig4Test

type TGateWay struct {
	ConnAddr string `json:"ConnAddr"`
}

type TArrGateWay []*TGateWay

var (
	_testCfgPath string
)

func GetTestCfgPath() (path string) {
	exepath := utls.GetExeFilePath()
	path = exepath + "/testconfig/"
	return
}

func init() {
	_testCfgPath = GetTestCfgPath()
	gameMsg.Init()
}
