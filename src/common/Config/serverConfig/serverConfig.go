package serverConfig

import "common/Define"

type TServerBaseConfig struct {
	ServerType       Define.ERouteId
	RedisConfig      *TRedisConfig
	MgoConfig        *TMgoConfig
	ExternalGWConfig *TGatewayConfig
	InnerGWConfig    *TGatewayConfig
	LoginConfig      *TLoginConfig
	GameConfig       *TGameConfig
}

var (
	GServerBaseConfig = &TServerBaseConfig{
		RedisConfig:      GRedisCfgProvider,
		MgoConfig:        GMgoCfgProvider,
		ExternalGWConfig: GExternalGWConfig,
		InnerGWConfig:    GInnerGWConfig,
		LoginConfig:      GLoginConfig,
		GameConfig:       GGameConfig,
	}
)
