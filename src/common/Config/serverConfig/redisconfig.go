package serverConfig

import (
	"common/Config"
	"fmt"
	"path/filepath"
	"strconv"
)

/*
	export from redisconfig.json by tool.
*/
type TRedisconfigBase struct {
	Id            int32  `json:"id"`
	DBIndex       int32  `json:"DBIndex"`
	Connaddr      string `json:"ConnAddr"`
	Shareconnaddr string `json:"ShareConnAddr"`
	Passwd        string `json:"Passwd"`
	Sharedbindex  int32  `json:"ShareDBIndex"`
	Pprofaddr     string `json:"PProfAddr"`
	Name          string
}

type TRedisconfigConfig struct {
	data *TRedisconfigBase
}

type tArrRedisconfig []*TRedisconfigBase

var (
	GRedisconfigConfig *TRedisconfigConfig = &TRedisconfigConfig{}
	cstRedisDef                            = "Redis"
)

func init() {
	//loadRedisConfig()
}

func loadRedisConfig() {
	var (
		redispath string
	)
	if len(SvrPath) == 0 {
		redispath = getserverpath()
	}
	redispath = filepath.Join(SvrPath, "redisconfig.json")
	Config.ParseJson2Cache(GRedisconfigConfig, &tArrRedisconfig{}, redispath)
}

func (this *TRedisconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrRedisconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Connaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("redisconfig Connaddr invalid, id: %v.", item.Id))
		}

		if len(item.Shareconnaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("redisconfig Shareconnaddr invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("redisconfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TRedisconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrRedisconfig)
	errlist = []string{}
	for _, item := range *cfg {
		item.Name = cstRedisDef + "_" + strconv.Itoa(int(item.Id))
		this.data = item
		break
	}
	return
}

func (this *TRedisconfigConfig) Get() *TRedisconfigBase {
	return this.data
}
