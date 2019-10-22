package serverConfig

type TLoginConfig struct {
	ListenAddr string
}

var (
	GLoginConfig = &TLoginConfig{
		ListenAddr: "0.0.0.0:51001",
	}
)
