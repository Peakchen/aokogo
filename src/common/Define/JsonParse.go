package Define

import (
	"encoding/json"
	"io/ioutil"
)

var (
	GJsonParseTool *TJsonParseTool
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
	GJsonParseTool = NewJsonParseTool()
}
