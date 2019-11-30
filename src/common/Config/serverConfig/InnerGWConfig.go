package serverConfig

import (
	"common/Config"
	"fmt"
)

/*
	export from InnerGWConfig.json by tool.
*/
type TInnergwconfigBase struct {
	Id          int32  `json:"id"`
	Connectaddr string `json:"ConnectAddr"`
	Listenaddr  string `json:"ListenAddr"`
	Zone        string `json:"Zone"`
	No          string `json:"No"`
	Pprofaddr   string `json:"PProfAddr"`
}

type TInnergwconfigConfig struct {
	data *TInnergwconfigBase
}

type tArrInnergwconfig []*TInnergwconfigBase

var (
	GInnergwconfigConfig *TInnergwconfigConfig = &TInnergwconfigConfig{}
)

func init() {
	Config.ParseJson2Cache(GInnergwconfigConfig, &tArrInnergwconfig{}, getserverpath()+"InnerGWConfig.json")
}

func (this *TInnergwconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrInnergwconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Listenaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("InnerGWConfig listeraddr invalid, id: %v.", item.Id))
		}

		if len(item.Connectaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("InnerGWConfig Connectaddr invalid, id: %v.", item.Id))
		}

		if len(item.Zone) == 0 {
			errlist = append(errlist, fmt.Sprintf("InnerGWConfig Zone invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("InnerGWConfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TInnergwconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrInnergwconfig)
	this.data = &TInnergwconfigBase{}
	for _, item := range *cfg {
		this.data = item
		break
	}
	return
}

func (this *TInnergwconfigConfig) Get() *TInnergwconfigBase {
	return this.data
}
