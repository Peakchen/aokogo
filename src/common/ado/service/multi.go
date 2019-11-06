package service

import "common/Config/serverConfig"

var (
	clusterProvider *TClusterDBProvider
)

func init() {
	clusterProvider = &TClusterDBProvider{}
}

func StartMultiDBProvider(Server string, RedisCfg *serverConfig.TRedisConfig, MgoCfg *serverConfig.TMgoConfig) {
	clusterProvider.Start(Server, RedisCfg, MgoCfg)
}
