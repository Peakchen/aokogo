package mgoDB

import (
	"github.com/globalsign/mgo"
	"common/Log"
	"time"
)

type TDBServer struct {
	sess *mgo.Session
	UserName string
	Passwd string
	ServiceHost string
}

func Create(Username, Passwd, host string)*TDBServer{
	dbsess := &TDBServer{}
	dbsess.UserName = Username
	dbsess.Passwd = Passwd
	dbsess.ServiceHost = host

	return dbsess
}

func (self *TDBServer) NewMgoServer(){
	MdialInfo := &mgo.DialInfo{
		Addrs: []string{self.ServiceHost},
		Username: self.UserName,
		Password: self.Passwd,
		Direct: false,
		Timeout: time.Second*3,
		PoolLimit: 4096,
		ReadTimeout: time.Second*5,
		WriteTimeout: time.Second*5,
	}

	session, err := mgo.DialWithInfo(MdialInfo)
	if err != nil {
		Log.Error("mgo dial err: %v.\n", err)
		return
	}
	
	session.SetMode(mgo.Monotonic,true)
	self.sess = session
	return
}

func (self *TDBServer) Stop(){
	if self.sess != nil {
		self.sess.Close()
	}
}

func (self *TDBServer) GetDB()*mgo.Session{
	if self.sess == nil {
		self.NewMgoServer()
	}

	return self.sess.Clone()
}

func (self *TDBServer) OnTimer2FlushDB(){
	reach := time.NewTicker(100*time.Millisecond)
	for {
		select{
		case <-reach.C:
			// todo:
			self.FlushDB()
		default:
			// nothing...

		}
	}
}

func (self *TDBServer) FlushDB(){
	
}

func (self *TDBServer) QueryOne(Model string, InParam interface{}, OutParam interface{}){
	
}

func (self *TDBServer) QuerySome(Model string, InParam interface{}, OutParam interface{}){

}

func (self *TDBServer) SaveOne(Model string, InParam interface{}, OutParam interface{}){

}