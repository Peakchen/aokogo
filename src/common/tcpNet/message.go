package tcpNet

import (
	"common/Define"
	"common/Log"
	"common/msgProto/MSG_Login"
	"common/msgProto/MSG_MainModule"
	"common/msgProto/MSG_Server"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
)

type TMessageProc struct {
	proc       reflect.Value
	paramTypes []reflect.Type
}

var (
	_MessageTab      map[uint32]*TMessageProc = map[uint32]*TMessageProc{}
	_specialLoginMsg                          = map[uint16]bool{}
)

func RegisterMessage(mainID, subID uint16, proc interface{}) {
	_cmd := EncodeCmd(mainID, subID)
	_, ok := _MessageTab[_cmd]
	if ok {
		return
	}

	cbref := reflect.TypeOf(proc)
	if cbref.Kind() != reflect.Func {
		Log.FmtPrintln("proc type not is func, but is: %v.", cbref.Kind())
		return
	}

	if cbref.NumIn() != 2 {
		Log.FmtPrintln("proc num input is not 2, but is: %v.", cbref.NumIn())
		return
	}

	if cbref.NumOut() != 2 {
		Log.FmtPrintln("proc num output is not 2, but is: %v.", cbref.NumOut())
		return
	}

	if cbref.Out(0) != reflect.TypeOf(bool(false)) {
		Log.FmtPrintln("proc num out 1 is not string, but is: %v.", cbref.Out(0))
		return
	}

	if cbref.Out(1).Name() != "error" {
		Log.FmtPrintln("proc num out 2 is not *proto.Message, but is: %v.", cbref.Out(1), reflect.TypeOf(error(nil)), errors.New("0"), fmt.Errorf("0"))
		return
	}

	paramtypes := []reflect.Type{}
	for i := 0; i < cbref.NumIn(); i++ {
		t := cbref.In(i)
		// if t.Kind() == reflect.String ||
		// 	t.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		// 	paramtypes = append(paramtypes, t)
		// }
		paramtypes = append(paramtypes, t)
	}

	_MessageTab[_cmd] = &TMessageProc{
		proc:       reflect.ValueOf(proc),
		paramTypes: paramtypes,
	}

	return
}

func GetMessageInfo(mainID, subID uint16) (proc *TMessageProc, finded bool) {
	_cmd := EncodeCmd(mainID, subID)
	proc, finded = _MessageTab[_cmd]
	return
}

func GetAllMessageIDs() (msgs []uint32) {
	msgs = []uint32{}
	for msgid, _ := range _MessageTab {
		msgs = append(msgs, uint32(msgid))
	}
	return
}

func RegisterMessageRet(session TcpSession, ServerType uint16) (succ bool, err error) {
	rsp := &MSG_Server.SC_ServerRegister_Rsp{}
	rsp.Ret = MSG_Server.ErrorCode_Success
	return session.SendMsg(ServerType,
		uint16(MSG_MainModule.MAINMSG_SERVER),
		uint16(MSG_Server.SUBMSG_SC_ServerRegister),
		rsp)
}

func SpecialLoginMsgFilter(main, sub uint16) (ok bool) {
	if main != uint16(MSG_MainModule.MAINMSG_LOGIN) {
		return
	}

	if sub == uint16(MSG_Login.SUBMSG_CS_UserRegister) ||
		sub == uint16(MSG_Login.SUBMSG_CS_Login) {
		ok = true
	}

	return
}

/*
@func: UnPackExternalMsg 解服务器外部消息（客户端，clientsession 注册消息）
@parma1: 连接对象 c *net.TCPConn
@param2: 解包对象 pack IMessagePack
*/
func UnPackExternalMsg(c *net.TCPConn, pack IMessagePack) (succ bool) {
	packLenBuf := make([]byte, EnMessage_NoDataLen)
	readn, err := io.ReadFull(c, packLenBuf)
	if err != nil || readn < EnMessage_NoDataLen {
		if err.Error() == "EOF" {
			succ = true
		} else {
			Log.FmtPrintln("read data fail, err: ", err, readn)
		}
		return
	}

	Log.FmtPrintln("identify is empty, read data: ", len(packLenBuf))
	packlen := binary.LittleEndian.Uint32(packLenBuf[EnMessage_DataPackLen:EnMessage_NoDataLen])
	if packlen > maxMessageSize {
		Log.FmtPrintln("error receiving packLen:", packlen)
		return
	}

	data := make([]byte, EnMessage_NoDataLen+packlen)
	readn, err = io.ReadFull(c, data[EnMessage_NoDataLen:])
	if err != nil || readn < int(packlen) {
		Log.FmtPrintln("error receiving msg, readn:", readn, "packLen:", packlen, "reason:", err)
		return
	}

	//todo: unpack message then read real date.
	copy(data[:EnMessage_NoDataLen], packLenBuf[:])
	_, err = pack.UnPackMsg4Client(data)
	if err != nil {
		Log.FmtPrintln("unpack action err: ", err)
		return
	}

	succ = true
	return
}

/*
@func: UnPackInnerMsg 解服务器内部消息（server 间客户端发来的请求或者其他rpc消息传递）
@parma1: 连接对象 c *net.TCPConn
@param2: 解包对象 pack IMessagePack
*/
func UnPackInnerMsg(c *net.TCPConn, pack IMessagePack) (succ bool) {
	packLenBuf := make([]byte, EnMessage_SvrNoDataLen)
	readn, err := io.ReadFull(c, packLenBuf)
	if err != nil || readn < EnMessage_SvrNoDataLen {
		if err.Error() == "EOF" {
			succ = true
		} else {
			Log.FmtPrintln("read data fail, err: ", err, readn)
		}
		return
	}

	Log.FmtPrintln("identify not empty, read data: ", len(packLenBuf))
	packlen := binary.LittleEndian.Uint32(packLenBuf[EnMessage_SvrDataPackLen:EnMessage_SvrNoDataLen])
	if packlen > maxMessageSize {
		Log.FmtPrintln("error receiving packLen:", packlen)
		return
	}

	data := make([]byte, EnMessage_SvrNoDataLen+packlen)
	readn, err = io.ReadFull(c, data[EnMessage_SvrNoDataLen:])
	if err != nil || readn < int(packlen) {
		Log.FmtPrintln("error receiving msg, readn:", readn, "packLen:", packlen, "reason:", err)
		return
	}

	//todo: unpack message then read real date.
	copy(data[:EnMessage_SvrNoDataLen], packLenBuf[:])
	_, err = pack.UnPackMsg4Svr(data)
	if err != nil {
		Log.FmtPrintln("unpack action err: ", err)
		return
	}
	succ = true
	return
}

/*
	内网关路由
*/
func innerMsgRouteAct(route, mainID uint16, data []byte) (succ bool) {
	var (
		session TcpSession
	)
	switch Define.ERouteId(route) {
	case Define.ERouteId_ER_ESG,
		Define.ERouteId_ER_ISG:
		if mainID == uint16(MSG_MainModule.MAINMSG_RPC) {
			//内网game rpc 调用
			Log.FmtPrintln("inner game rpc route.")
			session = GServer2ServerSession.GetSession(Define.ERouteId_ER_Game)
		} else {
			// 内网转发回复
			Log.FmtPrintln("inner respnse.")
			session = GServer2ServerSession.GetSession(Define.ERouteId_ER_ESG)
		}
	default:
		//内网转发路由请求
		Log.FmtPrintf("inner route requst message, route: %v.", route)
		session = GServer2ServerSession.GetSession(Define.ERouteId(route))
	}

	if session == nil {
		Log.FmtPrintf("can not find session from inner gateway, route: %v, mainID: %v.", route, mainID)
		return
	}

	if !session.Alive() {
		GServer2ServerSession.RemoveSession(session.GetIdentify())
	} else {
		succ = session.WriteMessage(data)
	}
	return
}

/*
	外网关路由
*/
func externalRouteAct(route uint16, obj *SvrTcpSession) (succ bool) {
	succ = true
	//客户端请求消息
	if Define.ERouteId(route) != Define.ERouteId_ER_ISG {
		Log.FmtPrintf("external request, route: %v, StrIdentify: %v.", route, obj.StrIdentify)
		GClient2ServerSession.AddSession(obj.RemoteAddr, obj)
		session := GServer2ServerSession.GetSession(Define.ERouteId_ER_ISG)
		if session == nil {
			Log.FmtPrintf("[request] can not find session route from external gateway, route: %v.", route)
			return
		}

		if !session.Alive() {
			GServer2ServerSession.RemoveSession(session.GetRemoteAddr())
		} else {
			out := make([]byte, EnMessage_SvrNoDataLen+int(obj.pack.GetDataLen()))
			err := obj.pack.PackAction(out)
			if err != nil {
				Log.FmtPrintln("[server] unpack action err: ", err)
				return
			}
			succ = session.WriteMessage(out)
		}
		return
	}

	//外网回复客户端消息
	Log.FmtPrintln("external response, StrIdentify: ", obj.pack.GetRemoteAddr(), len(obj.pack.GetRemoteAddr()))
	session := GClient2ServerSession.GetSessionByIdentify(obj.pack.GetRemoteAddr())
	if session == nil {
		Log.FmtPrintf("[response] can not find session route from external gateway, route: %v.", route)
		return
	}

	if !session.Alive() {
		GClient2ServerSession.RemoveSession(session.GetRemoteAddr())
	} else {
		out := make([]byte, EnMessage_NoDataLen+int(obj.pack.GetDataLen()))
		err := obj.pack.PackAction4Client(out)
		if err != nil {
			Log.FmtPrintln("[server] unpack action err: ", err)
			return
		}
		succ = session.WriteMessage(out)
	}
	return
}
