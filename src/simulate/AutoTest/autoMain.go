package AutoTest

import (
	"common/Log"
	"common/utls"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type TMsgDsc struct {
	SubId int32  `json:"subId"`
	Name  string `json:"name"`
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
	MainId int32 `json:"mainId"`
	Msg    *TMsg `json:"msg"`
}

type tArrConfig4Test []*TConfig4Test

var (
	_testCfgPath string
)

func getTestCfgPath() (path string) {
	exepath := utls.GetExeFilePath()
	path = exepath + "/testconfig/"
	return
}

func init() {
	_testCfgPath = getTestCfgPath()
}

func Run() {
	rd, err := ioutil.ReadDir(_testCfgPath)
	if err != nil {
		Log.Error("read test config path fail, err: ", err)
		return
	}

	for _, file := range rd {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if !strings.Contains(fileName, ".json") {
			continue
		}

		fileobj, err := ioutil.ReadFile(_testCfgPath + fileName)
		if err != nil {
			Log.Error("read config fail, info: ", fileName, err)
			return
		}

		// load module test msg
		data := &tArrConfig4Test{}
		err = json.Unmarshal(fileobj, &data)
		if err != nil {
			Log.Error("unmarshal config fail, info: ", fileName, err)
			return
		}

		// then run test...
	}
}
