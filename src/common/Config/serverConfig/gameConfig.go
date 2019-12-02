package serverConfig

import (
	"common/Config"
	"fmt"
)

/*
	export from gameConfig.json by tool.
*/
type TGameconfigBase struct {
	Id         int32  `json:"id"`
	No         string `json:"No"`
	Listenaddr string `json:"ListenAddr"`
	Zone       string `json:"Zone"`
	Pprofaddr  string `json:"PProfAddr"`
}

type TGameconfigConfig struct {
	data *TGameconfigBase
}

type tArrGameconfig []*TGameconfigBase

var (
	GGameconfigConfig *TGameconfigConfig = &TGameconfigConfig{}
)

func init() {
	Config.ParseJson2Cache(GGameconfigConfig, &tArrGameconfig{}, getserverpath()+"gameConfig.json")
}

func (this *TGameconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrGameconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Listenaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("gameConfig listeraddr invalid, id: %v.", item.Id))
		}

		if len(item.Zone) == 0 {
			errlist = append(errlist, fmt.Sprintf("gameConfig Zone invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("gameConfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TGameconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrGameconfig)
	this.data = &TGameconfigBase{}
	for _, item := range *cfg {
		this.data = item
		break
	}
	return
}

func (this *TGameconfigConfig) Get() *TGameconfigBase {
	return this.data
}
