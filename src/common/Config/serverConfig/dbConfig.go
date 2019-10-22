package serverConfig

type TRedisConfig struct {
	ConnAddr string
	DBIndex  int32
	Passwd   string
}

type TMgoConfig struct {
	UserName string
	Passwd   string
	Host     string
}

var GRedisCfgProvider *TRedisConfig = &TRedisConfig{
	ConnAddr: "0.0.0.0:6379",
	DBIndex:  1,
	Passwd:   "",
}

var GMgoCfgProvider *TMgoConfig = &TMgoConfig{
	Host:     "0.0.0.0:27017",
	UserName: "",
	Passwd:   "",
}
