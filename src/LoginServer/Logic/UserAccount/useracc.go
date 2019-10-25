package UserAccount

import (
	"LoginServer/dbo"
	"common/ado"
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

func RegisterUseAcc(acc *TUserAcc) (err error, exist bool) {
	acc.ModuleID = acc.UserName
	err = dbo.A_DBRead(acc)
	if err != nil {
		err = dbo.A_DBInsert(acc)
		if err != nil {
			return
		}
	} else {
		exist = true
	}
	return
}

func GetUserAcc(acc *TUserAcc) (exist bool) {
	acc.ModuleID = acc.UserName
	err := dbo.A_DBRead(acc)
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
