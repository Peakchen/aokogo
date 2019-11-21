package serverConfig

type TRedisConfig struct {
	//独有数据
	ConnAddr string
	DBIndex  int32
	Passwd   string

	//共享数据
	ShareConnAddr string
	ShareDBIndex  int32
	ShareDBPasswd string

	PProfAddr string
}

type TMgoConfig struct {
	//独有数据
	UserName string
	Passwd   string
	Host     string

	//共享数据
	ShareUserName string
	SharePasswd   string
	ShareHost     string

	PProfAddr string
}

var GRedisCfgProvider *TRedisConfig = &TRedisConfig{
	ConnAddr:  "0.0.0.0:6379",
	DBIndex:   1,
	Passwd:    "",
	PProfAddr: "11000",
}

var GMgoCfgProvider *TMgoConfig = &TMgoConfig{
	Host:      "0.0.0.0:27017",
	UserName:  "",
	Passwd:    "",
	PProfAddr: "11000",
}
