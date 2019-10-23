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

func (this *TUserAcc) CacheKey() string {
	return this.ModuleID
}

func (this *TUserAcc) MainModel() string {
	return cstAccMainModule
}

func (this *TUserAcc) SubModel() string {
	return cstAccSubModule
}

func FindUserAcc(acc *TUserAcc) (exist bool) {
	acc.ModuleID = bson.NewObjectId().Hex()
	err := dbo.A_DataGet(acc)
	if err != nil {

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
