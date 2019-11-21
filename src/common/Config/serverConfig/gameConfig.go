package serverConfig

type TGameConfig struct {
	No         string
	Zone       string
	ListenAddr string
	PProfAddr  string
}

var (
	GGameConfig = &TGameConfig{
		ListenAddr: "127.0.0.1:19000",
		PProfAddr:  "127.0.0.1:11001",
		No:         "1",
		Zone:       "Server",
	}
)
