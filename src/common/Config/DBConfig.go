package Config

import (

)

type TRedisConfig struct{
	ConnAddr string 
	DBIndex int32
	Passwd string
}

type TMgoConfig struct {
	UserName string
	Passwd string
	Host string
}

var GRedisCfgProvider *TRedisConfig = &TRedisConfig{}
var GMgoCfgProvider *TMgoConfig = &TMgoConfig{}
