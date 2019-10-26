package serverConfig

type TLoginConfig struct {
	No         string
	Zone       string
	ListenAddr string
}

var (
	GLoginConfig = &TLoginConfig{
		ListenAddr: "0.0.0.0:51001",
		No:         "1",
		Zone:       "Server",
	}
)
