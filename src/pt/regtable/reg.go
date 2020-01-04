package regtable

import (
	"common/tool"
	//"fmt"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

const (
	cstPortKey = `ZXCAD\Server_Client\loc`
	cstPortWordK = string("port")
)

/*
	purpose: check write port and exe path
*/
func CreateRegPortKey(srcPort int32) (nextPort int32) {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, cstPortKey, registry.ALL_ACCESS)
	if err != nil {
		tool.MyFmtPrint_Error("can not create register key, ")
		return
	}
	defer key.Close()

	//portStr := "port_" + strconv.Itoa(int(srcPort))
	err = key.SetStringValue(cstPortWordK, strconv.Itoa(int(srcPort)))
	if err != nil {
		tool.MyFmtPrint_Error("set int32 port fail, err: ", err)
		return
	}

	nextPort = srcPort

	err = key.SetStringValue("ExePath", tool.GetExeFilePath())
	if err != nil {
		tool.MyFmtPrint_Error("set int32 port fail, err: ", err)
		return
	}
	return
}

func GetRegPoryKey() (pv uint32){
	k, err := registry.OpenKey(registry.CURRENT_USER, cstPortKey, registry.ALL_ACCESS)
	if err != nil {
		tool.MyFmtPrint_Error(err)
		return
	}

	defer k.Close()
	value, _, err := k.GetStringValue(cstPortWordK)
	if err != nil {
		tool.MyFmtPrint_Error("get port key:", err)
		return
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		tool.MyFmtPrint_Error(err)
		return
	}

	pv = uint32(v)
	return 
}