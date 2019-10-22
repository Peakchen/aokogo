package DataProvider

import (
	"common/ado"
)

var (
	GDBProvider *ado.TDBProvider
)

func A_DataGet(DBKey string, Out interface{}){
	// check redis can get db data, if not exist, then from mogo.
	
}	

func A_DataSet(DBKey string, In interface{}){
	// check save data to redis cache or db persistence.

}