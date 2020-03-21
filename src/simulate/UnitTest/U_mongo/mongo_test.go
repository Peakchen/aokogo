package U_mongo

import (
	"testing"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"common/MgoConn"
	"common/Log"
)

func TestNormal(t *testing.T){
	var (
		Username string
		Passwd string
		Host string
	)
	mgoobj := MgoConn.NewMgoConn("test", Username, Passwd, Host)
	session, err := mgoobj.GetMgoSession()
	if err != nil {
		Log.Error(err)
		return
	}

	
}