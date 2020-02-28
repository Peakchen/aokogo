package serverConfig

import (
	"common/Config"
	"fmt"
	"path/filepath"
	"strconv"
)

/*
	export from InnerGWConfig.json by tool.
*/
type TInnergwconfigBase struct {
	Id          int32  `json:"id"`
	Connectaddr string `json:"ConnectAddr"`
	Listenaddr  string `json:"ListenAddr"`
	Zone        string `json:"Zone"`
	No          int32  `json:"No"`
	Pprofaddr   string `json:"PProfAddr"`
}

type TInnergwconfig struct {
	Id          int32
	Connectaddr string
	Listenaddr  string
	Zone        string
	No          string
	Pprofaddr   string
	Name        string
}

type TInnergwconfigConfig struct {
	data *TInnergwconfig
}

type tArrInnergwconfig []*TInnergwconfigBase

var (
	GInnergwconfigConfig *TInnergwconfigConfig = &TInnergwconfigConfig{}
	cstInnerGatewayDef                         = "InnerGateway"
)

func init() {
	//loadInnergwConfig()
}

func loadInnergwConfig() {
	var (
		InnerGWpath string
	)
	if len(SvrPath) == 0 {
		InnerGWpath = getserverpath()
	}
	InnerGWpath = filepath.Join(SvrPath, "InnerGWConfig.json")
	Config.ParseJson2Cache(GInnergwconfigConfig, &tArrInnergwconfig{}, InnerGWpath)
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
	this.data = &TInnergwconfig{}
	for _, item := range *cfg {
		num := strconv.Itoa(int(item.No))
		this.data = &TInnergwconfig{
			Id:          item.Id,
			No:          num,
			Connectaddr: item.Connectaddr,
			Listenaddr:  item.Listenaddr,
			Zone:        item.Zone,
			Pprofaddr:   item.Pprofaddr,
			Name:        cstInnerGatewayDef + "_" + strconv.Itoa(int(item.Id)),
		}
		break
	}
	return
}

func (this *TInnergwconfigConfig) Get() *TInnergwconfig {
	return this.data
}
