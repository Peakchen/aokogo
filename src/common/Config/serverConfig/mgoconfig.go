package serverConfig

import (
	"common/Config"
	"fmt"
)

/*
	export from mgoconfig.json by tool.
*/
type TMgoconfigBase struct {
	Id            int32  `json:"id"`
	Passwd        string `json:"Passwd"`
	Username      string `json:"UserName"`
	Shareusername string `json:"ShareUserName"`
	Host          string `json:"Host"`
	Sharehost     string `json:"ShareHost"`
	Sharepasswd   string `json:"SharePasswd"`
	Pprofaddr     string `json:"PProfAddr"`
}

type TMgoconfigConfig struct {
	data *TMgoconfigBase
}

type tArrMgoconfig []*TMgoconfigBase

var (
	GMgoconfigConfig *TMgoconfigConfig = &TMgoconfigConfig{}
)

func init() {
	Config.ParseJson2Cache(GMgoconfigConfig, &tArrMgoconfig{}, getserverpath()+"mgoconfig.json")
}

func (this *TMgoconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrMgoconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Host) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Host invalid, id: %v.", item.Id))
		}

		if len(item.Sharehost) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Sharehost invalid, id: %v.", item.Id))
		}

		if len(item.Sharehost) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Sharehost invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TMgoconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrMgoconfig)
	errlist = []string{}
	this.data = &TMgoconfigBase{}
	for _, item := range *cfg {
		this.data = item
		break
	}
	return
}

func (this *TMgoconfigConfig) Get() *TMgoconfigBase {
	return this.data
}
