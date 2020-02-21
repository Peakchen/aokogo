package serverConfig

import (
	"common/Config"
	"fmt"
)

/*
	export from ExternalGWConfig.json by tool.
*/
type TExternalgwconfigBase struct {
	Id         int32  `json:"id"`
	Listenaddr string `json:"ListenAddr"`
	Pprofaddr  string `json:"PProfAddr"`
}

type TExternalgwconfigConfig struct {
	data *TExternalgwconfigBase
}

type tArrExternalgwconfig []*TExternalgwconfigBase

var (
	GExternalgwconfigConfig *TExternalgwconfigConfig = &TExternalgwconfigConfig{}
)

func init() {
	loadExternalgwConfig()
}

func loadExternalgwConfig(){
	Config.ParseJson2Cache(GExternalgwconfigConfig, &tArrExternalgwconfig{}, getserverpath()+"ExternalGWConfig.json")
}

func (this *TExternalgwconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrExternalgwconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Listenaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("ExternalGWConfig listeraddr invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("ExternalGWConfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TExternalgwconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrExternalgwconfig)
	this.data = &TExternalgwconfigBase{}
	for _, item := range *cfg {
		this.data = item
		break
	}
	return
}

func (this *TExternalgwconfigConfig) Get() *TExternalgwconfigBase {
	return this.data
}
