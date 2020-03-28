package U_mongo

import (
	"common/Log"
	"common/MgoConn"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"testing"
)

func TestNormal(t *testing.T) {
	var (
		Username string
		Passwd   string
		Host     string
	)
	mgoobj := MgoConn.NewMgoConn("test", Username, Passwd, Host)
	session, err := mgoobj.GetMgoSession()
	if err != nil {
		Log.Error(err)
		return
	}

}
