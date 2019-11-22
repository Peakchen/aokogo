package serverConfig

type TGatewayConfig struct {
	ListenAddr string
	PProfAddr  string
}

type TGWInnerConfig struct {
	ListenAddr  string
	ConnectAddr string
	PProfAddr   string
	No          string
	Zone        string
}

var (
	GExternalGWConfig = &TGatewayConfig{
		ListenAddr: "0.0.0.0:51001",
		PProfAddr:  "127.0.0.1:12002",
	}

	GInnerGWConfig = &TGWInnerConfig{
		ListenAddr:  "0.0.0.0:19000",
		ConnectAddr: "0.0.0.0:51001",
		PProfAddr:   "127.0.0.1:12003",
		No:          "1",
		Zone:        "Server",
	}
)
