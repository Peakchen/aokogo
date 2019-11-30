package service

var (
	clusterProvider *TClusterDBProvider
)

func init() {
	clusterProvider = &TClusterDBProvider{}
}

func StartMultiDBProvider(Server string) {
	clusterProvider.Start(Server)
}
