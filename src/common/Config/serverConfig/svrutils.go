package serverConfig

import "common/utls"

func getserverpath() (path string) {
	exepath := utls.GetExeFilePath()
	path = exepath + "/serverconfig/"
	return
}
