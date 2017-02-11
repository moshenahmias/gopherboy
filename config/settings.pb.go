// Code generated by protoc-gen-go.
// source: settings.proto
// DO NOT EDIT!

/*
Package config is a generated protocol buffer package.

It is generated from these files:
	settings.proto

It has these top-level messages:
	Settings
*/
package config

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

type EJoypad int32

const (
	EJoypad_JoypadUp     EJoypad = 0
	EJoypad_JoypadDown   EJoypad = 1
	EJoypad_JoypadLeft   EJoypad = 2
	EJoypad_JoypadRight  EJoypad = 3
	EJoypad_JoypadA      EJoypad = 4
	EJoypad_JoypadB      EJoypad = 5
	EJoypad_JoypadSelect EJoypad = 6
	EJoypad_JoypadStart  EJoypad = 7
)

var EJoypad_name = map[int32]string{
	0: "JoypadUp",
	1: "JoypadDown",
	2: "JoypadLeft",
	3: "JoypadRight",
	4: "JoypadA",
	5: "JoypadB",
	6: "JoypadSelect",
	7: "JoypadStart",
}
var EJoypad_value = map[string]int32{
	"JoypadUp":     0,
	"JoypadDown":   1,
	"JoypadLeft":   2,
	"JoypadRight":  3,
	"JoypadA":      4,
	"JoypadB":      5,
	"JoypadSelect": 6,
	"JoypadStart":  7,
}

func (x EJoypad) String() string {
	return proto.EnumName(EJoypad_name, int32(x))
}
func (EJoypad) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Settings struct {
	JoypadMapping map[int32]EJoypad `protobuf:"bytes,1,rep,name=joypad_mapping,json=joypadMapping" json:"joypad_mapping,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value,enum=config.EJoypad"`
	SoundDevice   int32             `protobuf:"varint,2,opt,name=sound_device,json=soundDevice" json:"sound_device,omitempty"`
	Fps           uint32            `protobuf:"varint,3,opt,name=fps" json:"fps,omitempty"`
	Scale         uint32            `protobuf:"varint,4,opt,name=scale" json:"scale,omitempty"`
	Color_0       uint32            `protobuf:"varint,5,opt,name=color_0,json=color0" json:"color_0,omitempty"`
	Color_1       uint32            `protobuf:"varint,6,opt,name=color_1,json=color1" json:"color_1,omitempty"`
	Color_2       uint32            `protobuf:"varint,7,opt,name=color_2,json=color2" json:"color_2,omitempty"`
	Color_3       uint32            `protobuf:"varint,8,opt,name=color_3,json=color3" json:"color_3,omitempty"`
}

func (m *Settings) Reset()                    { *m = Settings{} }
func (m *Settings) String() string            { return proto.CompactTextString(m) }
func (*Settings) ProtoMessage()               {}
func (*Settings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Settings) GetJoypadMapping() map[int32]EJoypad {
	if m != nil {
		return m.JoypadMapping
	}
	return nil
}

func init() {
	proto.RegisterType((*Settings)(nil), "config.Settings")
	proto.RegisterEnum("config.EJoypad", EJoypad_name, EJoypad_value)
}

func init() { proto.RegisterFile("settings.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x91, 0xc1, 0x4e, 0xab, 0x40,
	0x14, 0x86, 0x2f, 0x50, 0xa0, 0x39, 0xb4, 0x74, 0x72, 0x72, 0x13, 0x27, 0xae, 0xaa, 0xc6, 0xa4,
	0x71, 0x41, 0x5a, 0xba, 0x31, 0xee, 0x34, 0x75, 0xd3, 0xe8, 0x42, 0x1a, 0xd7, 0x0d, 0xd2, 0x69,
	0xa5, 0x22, 0x43, 0x60, 0x5a, 0xd3, 0x07, 0xf0, 0x1d, 0x7d, 0x1c, 0x61, 0xa6, 0xcd, 0x90, 0xb8,
	0x3b, 0xff, 0xf7, 0xfd, 0x43, 0xce, 0x0c, 0xe0, 0x57, 0x4c, 0x88, 0x34, 0xdf, 0x54, 0x41, 0x51,
	0x72, 0xc1, 0xd1, 0x49, 0x78, 0xbe, 0x4e, 0x37, 0x97, 0x3f, 0x26, 0x74, 0x17, 0x47, 0x85, 0x73,
	0xf0, 0xb7, 0xfc, 0x50, 0xc4, 0xab, 0xe5, 0x67, 0x5c, 0x14, 0x35, 0xa2, 0xc6, 0xd0, 0x1a, 0x79,
	0xe1, 0x55, 0xa0, 0xda, 0xc1, 0xa9, 0x19, 0xcc, 0x65, 0xed, 0x59, 0xb5, 0x1e, 0x73, 0x51, 0x1e,
	0xa2, 0xfe, 0xb6, 0xcd, 0xf0, 0x02, 0x7a, 0x15, 0xdf, 0xe5, 0xab, 0xe5, 0x8a, 0xed, 0xd3, 0x84,
	0x51, 0x73, 0x68, 0x8c, 0xec, 0xc8, 0x93, 0x6c, 0x26, 0x11, 0x12, 0xb0, 0xd6, 0x45, 0x45, 0xad,
	0xda, 0xf4, 0xa3, 0x66, 0xc4, 0xff, 0x60, 0x57, 0x49, 0x9c, 0x31, 0xda, 0x91, 0x4c, 0x05, 0x3c,
	0x03, 0x37, 0xe1, 0x19, 0x2f, 0x97, 0x63, 0x6a, 0x4b, 0xee, 0xc8, 0x38, 0xd6, 0x62, 0x42, 0x9d,
	0x96, 0x98, 0x68, 0x11, 0x52, 0xb7, 0x25, 0x42, 0x2d, 0xa6, 0xb4, 0xdb, 0x12, 0xd3, 0xf3, 0x17,
	0xc0, 0xbf, 0x77, 0x6a, 0x36, 0xfc, 0x60, 0x87, 0xfa, 0x15, 0x9a, 0xdd, 0x9b, 0x11, 0xaf, 0xc1,
	0xde, 0xc7, 0xd9, 0x4e, 0xdd, 0xc7, 0x0f, 0x07, 0xa7, 0x97, 0x61, 0xea, 0x74, 0xa4, 0xec, 0x9d,
	0x79, 0x6b, 0xdc, 0x7c, 0x1b, 0xe0, 0x1e, 0x31, 0xf6, 0xa0, 0xab, 0xa6, 0xd7, 0x82, 0xfc, 0x43,
	0x1f, 0x40, 0xa5, 0x19, 0xff, 0xca, 0x89, 0xa1, 0xf3, 0x13, 0x5b, 0x0b, 0x62, 0xe2, 0x00, 0xbc,
	0xe3, 0xe7, 0xd2, 0xcd, 0xbb, 0x20, 0x16, 0x7a, 0xe0, 0x2a, 0x70, 0x4f, 0x3a, 0x3a, 0x3c, 0x10,
	0xbb, 0xde, 0xb0, 0xa7, 0xc2, 0x82, 0x65, 0x2c, 0x11, 0xc4, 0xd1, 0x87, 0x17, 0x22, 0x2e, 0x05,
	0x71, 0xdf, 0x1c, 0xf9, 0xc7, 0xa7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x39, 0x46, 0x10, 0x16,
	0x03, 0x02, 0x00, 0x00,
}
