package serverConfig

func LoadSvrAllConfig(){
	loadExternalgwConfig()
	loadGameConfig()
	loadInnergwConfig()
	loadLoginConfig()
	loadMgoConfig()
	loadNetFilterConfig()
	loadRedisConfig()
	
}