package UserAccount

import (
	"LoginServer/dbo"
	"common/ado"

	"github.com/globalsign/mgo/bson"
)

/*
	User Account.
	foucs: username and passwd is made md5 marshal.
*/
type TUserAcc struct {
	ado.IDBModule

	UserName   string
	Passwd     string
	DeviceNo   string // device number
	DeviceType string //device type (ios,androd)
}

func (this *TUserAcc) Identify() string {
	return this.StrIdentify
}

func (this *TUserAcc) MainModel() string {
	return cstAccMainModule
}

func (this *TUserAcc) SubModel() string {
	return cstAccSubModule
}

func RegisterUseAcc(acc *TUserAcc) (err error, exist bool) {
	err, exist = GetUserAcc(acc)
	if !exist {
		acc.StrIdentify = bson.NewObjectId().Hex()
		err = dbo.A_DBInsert(acc)
		if err != nil {
			return
		}
	}
	return
}

/*
	find user account by user name.
*/
func GetUserAcc(acc *TUserAcc) (err error, exist bool) {
	acc.StrIdentify = acc.UserName
	err, exist = dbo.A_DBReadAcc(acc)
	if err == nil {
		exist = true
	}
	return
}

/*
	find user account by user identify.
*/
func LoadUserAcc(Identify string) (data *TUserAcc, err error, exist bool) {
	data = &TUserAcc{}
	data.StrIdentify = Identify
	err, exist = dbo.A_DBRead(data)
	return
}
