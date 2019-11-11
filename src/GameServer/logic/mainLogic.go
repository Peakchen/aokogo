package logic

import (
	"common/Log"
	"reflect"
)

// after player login, need getting ready.
func EnterGameReady(identify string) {
	params := []reflect.Value{reflect.ValueOf(identify)}
	for module, obj := range GEnterReadyModule {
		enter := reflect.ValueOf(obj).MethodByName("EnterReady")
		if enter.IsNil() || !enter.IsValid() {
			Log.ErrorIDCard(identify, "can not find EnterReady method, module: ", module)
			return
		}

		enter.Call(params)
	}
}

//before leave, get ready.
func LeaveGameReady(identify string) {
	params := []reflect.Value{reflect.ValueOf(identify)}
	for module, obj := range GLeaveReadyModule {
		enter := reflect.ValueOf(obj).MethodByName("LeaveReady")
		if enter.IsNil() || !enter.IsValid() {
			Log.ErrorIDCard(identify, "can not find LeaveReady method, module: ", module)
			return
		}

		enter.Call(params)
	}
}

//broken link reconnct enter game ready.
func ReconnectEnterReady(identify string) {
	params := []reflect.Value{reflect.ValueOf(identify)}
	for module, obj := range GReconnReadyModule {
		enter := reflect.ValueOf(obj).MethodByName("ReconnectReady")
		if enter.IsNil() || !enter.IsValid() {
			Log.ErrorIDCard(identify, "can not find ReconnectReady method, module: ", module)
			return
		}

		enter.Call(params)
	}
}
