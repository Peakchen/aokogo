// Code generated by protoc-gen-go.
// source: MSG_MainModule.proto
// DO NOT EDIT!

/*
Package MSG_MainModule is a generated protocol buffer package.

It is generated from these files:
	MSG_MainModule.proto

It has these top-level messages:
*/
package MSG_MainModule

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// add by stefanchen
// 主模块消息结构
type MAINMSG int32

const (
	MAINMSG_Begin  MAINMSG = 0
	MAINMSG_SERVER MAINMSG = 1
	MAINMSG_LOGIN  MAINMSG = 2
)

var MAINMSG_name = map[int32]string{
	0: "Begin",
	1: "SERVER",
	2: "LOGIN",
}
var MAINMSG_value = map[string]int32{
	"Begin":  0,
	"SERVER": 1,
	"LOGIN":  2,
}

func (x MAINMSG) String() string {
	return proto.EnumName(MAINMSG_name, int32(x))
}
func (MAINMSG) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterEnum("MSG_MainModule.MAINMSG", MAINMSG_name, MAINMSG_value)
}

func init() { proto.RegisterFile("MSG_MainModule.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 98 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xf1, 0x0d, 0x76, 0x8f,
	0xf7, 0x4d, 0xcc, 0xcc, 0xf3, 0xcd, 0x4f, 0x29, 0xcd, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xe2, 0x43, 0x15, 0xd5, 0xd2, 0xe6, 0x62, 0xf7, 0x75, 0xf4, 0xf4, 0xf3, 0x0d, 0x76, 0x17,
	0xe2, 0xe4, 0x62, 0x75, 0x4a, 0x4d, 0xcf, 0xcc, 0x13, 0x60, 0x10, 0xe2, 0xe2, 0x62, 0x0b, 0x76,
	0x0d, 0x0a, 0x73, 0x0d, 0x12, 0x60, 0x04, 0x09, 0xfb, 0xf8, 0xbb, 0x7b, 0xfa, 0x09, 0x30, 0x25,
	0xb1, 0x81, 0xcd, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x7c, 0x1b, 0x1e, 0x5b, 0x00,
	0x00, 0x00,
}