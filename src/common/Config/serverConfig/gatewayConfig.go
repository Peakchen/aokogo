package serverConfig

type TGatewayConfig struct {
	ListenAddr string
}

type TGWInnerConfig struct {
	ListenAddr  string
	ConnectAddr string
	No          string
	Zone        string
}

var (
	GExternalGWConfig = &TGatewayConfig{
		ListenAddr: "0.0.0.0:51001",
	}

	GInnerGWConfig = &TGWInnerConfig{
		ListenAddr:  "0.0.0.0:19000",
		ConnectAddr: "0.0.0.0:51001",
		No:          "1",
		Zone:        "Server",
	}
)
