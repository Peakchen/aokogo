package MgoService

import (
	"github.com/globalsign/mgo"
	"common/Log"
	"time"
	"github.com/globalsign/mgo/bson"
	"common/RedisService"
)

type MgoConn struct {
	session *mgo.Session
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
	self.session = session
	return
}

func (self *MgoConn) Stop(){
	if self.session != nil {
		self.session.Close()
	}
}

func (self *MgoConn) GetDB()*mgo.Session{
	if self.session == nil {
		self.NewMgoConn()
	}

	return self.session.Clone()
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

func (self *MgoConn) QueryOne(OutParam IDBCache)(err error){
	collection := self.session.DB(OutParam.MainModel()).C(OutParam.SubModel())
	err = collection.Find(bson.M{"_id": OutParam.CacheKey()}).One(&OutParam)
	if err != nil {
		err = fmt.Errorf("CacheKey: %v, MainModel: %v, SubModel: %v, err: %v.\n", OutParam.CacheKey(), OutParam.MainModel(), OutParam.SubModel(), err)
		Log.Error("[QueryOne] err: %v.\n", err)
		return
	}
	return
}

func (self *MgoConn) QuerySome(OutParam IDBCache)(err error){
	collection := self.session.DB(OutParam.MainModel()).C(OutParam.SubModel())
	err = collection.Find(bson.M{"_id": OutParam.CacheKey()}).All(&OutParam)
	if err != nil {
		err = fmt.Errorf("CacheKey: %v, MainModel: %v, SubModel: %v, err: %v.\n", OutParam.CacheKey(), OutParam.MainModel(), OutParam.SubModel(), err)
		Log.Error("[QuerySome] err: %v.\n", err)
		return
	}
	return
}

func (self *MgoConn) SaveOne(InParam IDBCache)(err error){
	collection := self.session.DB(InParam.MainModel()).C(InParam.SubModel())
	operAction := bson.M{"$set": bson.M{InParam.SubModel(): InParam}}
	err = collection.Update(bson.M{"_id": InParam.CacheKey()}, operAction)
	if err != nil {
		err = fmt.Errorf("CacheKey: %v, MainModel: %v, SubModel: %v, err: %v.\n", InParam.CacheKey(), InParam.MainModel(), InParam.SubModel(), err)
		Log.Error("[SaveOne] err: %v.\n", err)
		return
	}

	return
}