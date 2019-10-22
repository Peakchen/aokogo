package UserAccount

import "common/ado"

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

/*

 */
func Get(module ado.IDBModule) (err error) {

	return
}

/*

 */
func Insert(module ado.IDBModule) (err error) {

	return
}

/*
 */
func Update(module ado.IDBModule) (err error) {

	return
}

/*
 */
func Delete(module ado.IDBModule) (err error) {

}
