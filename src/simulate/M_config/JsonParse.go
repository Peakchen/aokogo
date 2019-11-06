package M_config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	_JsonParseTool *TJsonParseTool
)

type TJsonParseTool struct {
}

func (this *TJsonParseTool) Parse(jsonName string, data interface{}) (err error) {
	filedata, err := ioutil.ReadFile(jsonName)
	if err != nil {
		return
	}

	err = json.Unmarshal(filedata, data)
	return
}

func NewJsonParseTool() *TJsonParseTool {
	return &TJsonParseTool{}
}

func init() {
	_JsonParseTool = NewJsonParseTool()
}
