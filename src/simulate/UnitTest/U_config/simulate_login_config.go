package U_config

import (
	"common/Config"
	"common/define"
	"common/utls"
	"fmt"
)

/*
	simulate test register and login.
*/
type TSimulateLoginBase struct {
	Username    string               `json:"username"`
	Passwd      string               `json:"passwd"`
	Register    int32                `json:"register"`
	Login       int32                `json:"login"`
	List        define.Int32Array    `json:"list"`
	List2D      define.Int32Array2D  `json:"list2D"`
	Property    define.Property      `json:"property"`
	PropertyArr define.PropertyArray `json:"propertylist"`
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
	data []*TSimulateLoginBase
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
	errlist = []string{}
	for idx, item := range *cfg {
		if len(item.Username) == 0 {
			errlist = append(errlist, fmt.Sprintf("user name invalid, idx: %v.", idx))
		}

		if len(item.Passwd) == 0 {
			errlist = append(errlist, fmt.Sprintf("user Passwd invalid, idx: %v.", idx))
		}
	}
	return
}

func (this *TSimulateLoginConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrSimulateLogin)
	this.data = []*TSimulateLoginBase{}
	for _, item := range *cfg {
		this.data = append(this.data, item)
	}
	return
}

func (this *TSimulateLoginConfig) Get() (data []*TSimulateLoginBase) {
	data = this.data
	return
}

func (this *TSimulateLoginConfig) GetItem(idx int32) (data *TSimulateLoginBase) {
	data = this.data[idx]
	return
}
