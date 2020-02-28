package serverConfig

import (
	"common/Config"
	"fmt"
	"path/filepath"
	"strconv"
)

/*
	export from gameConfig.json by tool.
*/
type TGameconfigBase struct {
	Id         int32  `json:"id"`
	No         int32  `json:"No"`
	Listenaddr string `json:"ListenAddr"`
	Zone       string `json:"Zone"`
	Pprofaddr  string `json:"PProfAddr"`
}

type TGameconfig struct {
	Id         int32
	No         string
	Listenaddr string
	Zone       string
	Pprofaddr  string
	Name       string
}

type TGameconfigConfig struct {
	data *TGameconfig
}

type tArrGameconfig []*TGameconfigBase

var (
	GGameconfigConfig *TGameconfigConfig = &TGameconfigConfig{}
	cstGameDef                           = "Game"
)

func init() {
	//loadGameConfig()
}

func loadGameConfig() {
	var (
		gameGWpath string
	)
	if len(SvrPath) == 0 {
		gameGWpath = getserverpath()
	}
	gameGWpath = filepath.Join(SvrPath, "gameConfig.json")
	Config.ParseJson2Cache(GGameconfigConfig, &tArrGameconfig{}, gameGWpath)
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
	this.data = &TGameconfig{}
	for _, item := range *cfg {
		num := strconv.Itoa(int(item.No))
		this.data = &TGameconfig{
			Id:         item.Id,
			No:         num,
			Listenaddr: item.Listenaddr,
			Zone:       item.Zone,
			Pprofaddr:  item.Pprofaddr,
			Name:       cstGameDef + "_" + strconv.Itoa(int(item.Id)),
		}
		break
	}
	return
}

func (this *TGameconfigConfig) Get() *TGameconfig {
	return this.data
}
