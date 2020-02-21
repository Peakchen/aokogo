package serverConfig

import (
	"common/Config"
	"fmt"
	"strconv"
)

/*
	export from LoginConfig.json by tool.
*/
type TLoginconfigBase struct {
	Id         int32  `json:"id"`
	No         int32  `json:"No"`
	Listenaddr string `json:"ListenAddr"`
	Zone       string `json:"Zone"`
	Pprofaddr  string `json:"PProfAddr"`
}

type TLoginconfig struct {
	Id         int32
	No         string
	Listenaddr string
	Zone       string
	Pprofaddr  string
}

type TLoginconfigConfig struct {
	data *TLoginconfig
}

type tArrLoginconfig []*TLoginconfigBase

var (
	GLoginconfigConfig *TLoginconfigConfig = &TLoginconfigConfig{}
)

func init() {
	loadLoginConfig()
}

func loadLoginConfig(){
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
	this.data = &TLoginconfig{}
	for _, item := range *cfg {
		num := strconv.Itoa(int(item.No))
		this.data = &TLoginconfig{
			Id:         item.Id,
			No:         num,
			Listenaddr: item.Listenaddr,
			Zone:       item.Zone,
			Pprofaddr:  item.Pprofaddr,
		}
		break
	}
	return
}

func (this *TLoginconfigConfig) Get() *TLoginconfig {
	return this.data
}
