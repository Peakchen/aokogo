package MgoService

import (
	"github.com/globalsign/mgo"
	"common/Log"
	"time"
	"github.com/globalsign/mgo/bson"
)

type MgoConn struct {
	sess *mgo.Session
	UserName string
	Passwd string
	ServiceHost string
}

func NewMgoConn(Username, Passwd, Host string)*MgoConn{
	dbsess := &MgoConn{}
	dbsess.UserName = Username
	dbsess.Passwd = Passwd
	dbsess.ServiceHost = Host

	dbsess.NewDial()
	return dbsess
}

func (self *MgoConn) NewDial(){
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

func (self *MgoConn) Stop(){
	if self.sess != nil {
		self.sess.Close()
	}
}

func (self *MgoConn) GetDB()*mgo.Session{
	if self.sess == nil {
		self.NewMgoConn()
	}

	return self.sess.Clone()
}

func (self *MgoConn) OnTimer2FlushDB(){
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

func (self *MgoConn) FlushDB(){
	
}

func MakeDBModel(Identify, MainModel, SubModel string)string {
	return MainModel+"."+SubModel+"."+Identify
}

func (self *MgoConn) QueryOne(Identify, MainModel, SubModel string, OutParam interface{}){
	//DBModel := MakeDBModel(Identify, MainModel, SubModel)
	collection := self.sess.DB(MainModel).C(SubModel)
	err := collection.Find(bson.M{"_id": Identify}).One(&OutParam)
	if err != nil {
		Log.Error("[QueryOne] Identify: %v, MainModel: %v, SubModel: %v.\n", Identify, MainModel, SubModel)
		return
	}

}

func (self *MgoConn) QuerySome(Identify, MainModel, SubModel string, OutParam interface{}){
	collection := self.sess.DB(MainModel).C(SubModel)
	err := collection.Find(bson.M{"_id": Identify}).All(&OutParam)
	if err != nil {
		Log.Error("[QuerySome] Identify: %v, MainModel: %v, SubModel: %v.\n", Identify, MainModel, SubModel)
		return
	}
}

func (self *MgoConn) SaveOne(Identify, MainModel, SubModel string,  InParam interface{}){
	collection := self.sess.DB(MainModel).C(SubModel)
	err := collection.Update(bson.M{"_id": Identify}, &InParam)
	if err != nil {
		Log.Error("[SaveOne] Identify: %v, MainModel: %v, SubModel: %v.\n", Identify, MainModel, SubModel)
		return
	}
}