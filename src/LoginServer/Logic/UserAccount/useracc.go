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

func GetUserAcc(acc *TUserAcc) (err error, exist bool) {
	acc.StrIdentify = acc.UserName
	err, exist = dbo.A_DBReadAcc(acc)
	if err == nil {
		exist = true
	}
	return
}

/*

 */
func (this *TUserAcc) Get(module ado.IDBModule) (err error) {

	return
}

/*

 */
func (this *TUserAcc) Insert(module ado.IDBModule) (err error) {

	return
}

/*
 */
func (this *TUserAcc) Update(module ado.IDBModule) (err error) {

	return
}

/*
 */
func (this *TUserAcc) Delete(module ado.IDBModule) (err error) {

	return
}
