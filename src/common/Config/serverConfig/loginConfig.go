package serverConfig

import (
	"common/Config"
	"fmt"
)

/*
	export from LoginConfig.json by tool.
*/
type TLoginconfigBase struct {
	Id         int32  `json:"id"`
	No         string `json:"No"`
	Listenaddr string `json:"ListenAddr"`
	Zone       string `json:"Zone"`
	Pprofaddr  string `json:"PProfAddr"`
}

type TLoginconfigConfig struct {
	data *TLoginconfigBase
}

type tArrLoginconfig []*TLoginconfigBase

var (
	GLoginconfigConfig *TLoginconfigConfig = &TLoginconfigConfig{}
)

func init() {
	Config.ParseJson2Cache(GLoginconfigConfig, &tArrLoginconfig{}, getserverpath()+"LoginConfig.json")
}

func (this *TLoginconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrLoginconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Listenaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("LoginConfig listeraddr invalid, id: %v.", item.Id))
		}

		if len(item.Zone) == 0 {
			errlist = append(errlist, fmt.Sprintf("LoginConfig Zone invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("LoginConfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TLoginconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrLoginconfig)
	this.data = &TLoginconfigBase{}
	for _, item := range *cfg {
		this.data = item
		break
	}
	return
}

func (this *TLoginconfigConfig) Get() *TLoginconfigBase {
	return this.data
}
