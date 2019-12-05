package U_config

import (
	"common/Config"
	"common/Define"
	"common/Log"
	"common/utls"
)

/*
	simulate test register and login.
*/
type TSimulateLoginBase struct {
	Username    string               `json:"username"`
	Passwd      string               `json:"passwd"`
	Register    int32                `json:"register"`
	Login       int32                `json:"login"`
	List        Define.Int32Array    `json:"list"`
	List2D      Define.Int32Array2D  `json:"list2D"`
	Property    Define.Property      `json:"property"`
	PropertyArr Define.PropertyArray `json:"propertylist"`
}

const (
	CstRegister_No  = int32(0)
	CstRegister_Yes = int32(1)
)

const (
	CstLogin_No  = int32(0)
	CstLogin_Yes = int32(1)
)

type TSimulateLoginConfig struct {
	data map[string]*TSimulateLoginBase
}

type tArrSimulateLogin []*TSimulateLoginBase

var (
	GloginConfig *TSimulateLoginConfig = &TSimulateLoginConfig{}
)

func getloginfile() (realfilename string) {
	exepath := utls.GetExeFilePath()
	realfilename = exepath + "/dataconfig/simulate_login.json"
	return
}

func init() {
	Config.ParseJson2Cache(GloginConfig, &tArrSimulateLogin{}, getloginfile())
}

func (this *TSimulateLoginConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrSimulateLogin)
	this.data = map[string]*TSimulateLoginBase{}
	for _, item := range *cfg {
		Log.FmtPrintln("ComfireAct act: ", item.Username, item.Passwd)
		this.data[item.Username] = item
	}
	return
}

func (this *TSimulateLoginConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrSimulateLogin)
	for _, item := range *cfg {
		Log.FmtPrintln("DataRWAct act: ", item.Username, item.Passwd)
		this.data[item.Username] = item
	}
	return
}

func (this *TSimulateLoginConfig) Get() (data map[string]*TSimulateLoginBase) {
	data = this.data
	return
}

func (this *TSimulateLoginConfig) GetItem(name string) (data *TSimulateLoginBase, exist bool) {
	data, exist = this.data[name]
	return
}
