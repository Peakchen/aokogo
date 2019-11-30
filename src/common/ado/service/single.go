package service

var (
	GDBProvider *TDBProvider
)

func NewDBProvider() {
	GDBProvider = &TDBProvider{}
}

func Run(server string) {
	GDBProvider.StartDBService(server)
}

func init() {
	NewDBProvider()
}
