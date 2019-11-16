package serverConfig

type TGameConfig struct {
	No         string
	Zone       string
	ListenAddr string
}

var (
	GGameConfig = &TGameConfig{
		ListenAddr: "127.0.0.1:19000",
		No:         "1",
		Zone:       "Server",
	}
)
