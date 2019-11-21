package serverConfig

type TLoginConfig struct {
	No         string
	Zone       string
	ListenAddr string
	PProfAddr  string
}

var (
	GLoginConfig = &TLoginConfig{
		ListenAddr: "0.0.0.0:19000",
		PProfAddr:  "127.0.0.1:11004",
		No:         "1",
		Zone:       "Server",
	}
)
