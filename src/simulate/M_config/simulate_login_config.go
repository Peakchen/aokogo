package M_config

import (
	"common/Define"
	"common/Log"
	"common/utls"
)

/*
	simulate test register and login.
*/
type TSimulateLogin struct {
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

type tArrSimulateLogin []*TSimulateLogin

var (
	Gloginconfig *tArrSimulateLogin = &tArrSimulateLogin{}
)

func getloginfile() (realfilename string) {
	exepath := utls.GetExeFilePath()
	realfilename = exepath + "/dataconfig/simulate_login.json"
	return
}

func init() {
	err := Define.GJsonParseTool.Parse(getloginfile(), Gloginconfig)
	if err != nil {
		Log.FmtPrintln("parse json fail, err: ", err)
		return
	}
	Log.FmtPrintln("login file: ", *Gloginconfig)
}
