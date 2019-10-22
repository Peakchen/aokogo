package serverConfig

type TGatewayConfig struct {
	ListenAddr string
}

var (
	GExternalGWConfig = &TGatewayConfig{
		ListenAddr: "0.0.0.0:51001",
	}

	GInnerGWConfig = &TGatewayConfig{
		ListenAddr: "0.0.0.0:19000",
	}
)
