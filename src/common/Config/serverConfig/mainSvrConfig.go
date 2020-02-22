package serverConfig

var (
	SvrPath string
)

func LoadSvrAllConfig(CfgPath string) {
	SvrPath = CfgPath
	loadExternalgwConfig()
	loadGameConfig()
	loadInnergwConfig()
	loadLoginConfig()
	loadMgoConfig()
	loadNetFilterConfig()
	loadRedisConfig()

}
