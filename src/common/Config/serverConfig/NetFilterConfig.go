package serverConfig

import (
	"common/Config"
	"fmt"
	"path/filepath"
)

/*
	export from NetFilter.json by tool.
*/
type TNetFilterBase struct {
	Id    int32  `json:"id"`
	White string `json:"white"`
	Black string `json:"black"`
}

type TNetFilter struct {
	Id    int32
	White string
	Black string
}

type TNetFilterConfig struct {
	data []*TNetFilter
}

type tArrNetFilter []*TNetFilterBase

var (
	GNetFilterConfig *TNetFilterConfig = &TNetFilterConfig{}
)

func init() {
	//loadNetFilterConfig()
}

func loadNetFilterConfig() {
	var (
		NetFilterpath string
	)
	if len(SvrPath) == 0 {
		NetFilterpath = getserverpath()
	}
	NetFilterpath = filepath.Join(SvrPath, "NetFilter.json")
	Config.ParseJson2Cache(GNetFilterConfig, &tArrNetFilter{}, NetFilterpath)
}

func (this *TNetFilterConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrNetFilter)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.White) == 0 {
			errlist = append(errlist, fmt.Sprintf("NetFilter White invalid, id: %v.", item.Id))
		}

		if len(item.Black) == 0 {
			errlist = append(errlist, fmt.Sprintf("NetFilter Black invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TNetFilterConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrNetFilter)
	this.data = []*TNetFilter{}
	for _, item := range *cfg {
		this.data = append(this.data, &TNetFilter{
			Id:    item.Id,
			White: item.White,
			Black: item.Black,
		})
	}
	return
}

func (this *TNetFilterConfig) Get() []*TNetFilter {
	return this.data
}
