package serverConfig

type TGameConfig struct {
	ListenAddr string
}

var (
	GGameConfig = &TGameConfig{
		ListenAddr: "127.0.0.1:19000",
	}
)
