package server

import (
	"common/ado/service"
)

/*
	run db server.
*/
func StartDBServer() {
	server := "sever1"
	service.StartMultiDBProvider(server)
}
