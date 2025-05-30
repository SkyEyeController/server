// Code generated by protoc-gen-go. DO NOT EDIT.
// source: network.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 定义 IDU 消息
type IDU struct {
	FromAddress          string   `protobuf:"bytes,1,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ToAddress            string   `protobuf:"bytes,2,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"`
	ContentType          string   `protobuf:"bytes,3,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Content              string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Timestamp            string   `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Proof                string   `protobuf:"bytes,6,opt,name=proof,proto3" json:"proof,omitempty"`
	AdditionalInfo       string   `protobuf:"bytes,7,opt,name=additional_info,json=additionalInfo,proto3" json:"additional_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDU) Reset()         { *m = IDU{} }
func (m *IDU) String() string { return proto.CompactTextString(m) }
func (*IDU) ProtoMessage()    {}
func (*IDU) Descriptor() ([]byte, []int) {
	return fileDescriptor_8571034d60397816, []int{0}
}

func (m *IDU) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDU.Unmarshal(m, b)
}
func (m *IDU) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDU.Marshal(b, m, deterministic)
}
func (m *IDU) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDU.Merge(m, src)
}
func (m *IDU) XXX_Size() int {
	return xxx_messageInfo_IDU.Size(m)
}
func (m *IDU) XXX_DiscardUnknown() {
	xxx_messageInfo_IDU.DiscardUnknown(m)
}

var xxx_messageInfo_IDU proto.InternalMessageInfo

func (m *IDU) GetFromAddress() string {
	if m != nil {
		return m.FromAddress
	}
	return ""
}

func (m *IDU) GetToAddress() string {
	if m != nil {
		return m.ToAddress
	}
	return ""
}

func (m *IDU) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *IDU) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *IDU) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *IDU) GetProof() string {
	if m != nil {
		return m.Proof
	}
	return ""
}

func (m *IDU) GetAdditionalInfo() string {
	if m != nil {
		return m.AdditionalInfo
	}
	return ""
}

// 定义 RelayRequest 消息类型，包含一个 IDU
type RelayRequest struct {
	Idu                  *IDU     `protobuf:"bytes,1,opt,name=idu,proto3" json:"idu,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelayRequest) Reset()         { *m = RelayRequest{} }
func (m *RelayRequest) String() string { return proto.CompactTextString(m) }
func (*RelayRequest) ProtoMessage()    {}
func (*RelayRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8571034d60397816, []int{1}
}

func (m *RelayRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelayRequest.Unmarshal(m, b)
}
func (m *RelayRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelayRequest.Marshal(b, m, deterministic)
}
func (m *RelayRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayRequest.Merge(m, src)
}
func (m *RelayRequest) XXX_Size() int {
	return xxx_messageInfo_RelayRequest.Size(m)
}
func (m *RelayRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RelayRequest proto.InternalMessageInfo

func (m *RelayRequest) GetIdu() *IDU {
	if m != nil {
		return m.Idu
	}
	return nil
}

// 定义 RelayResponse 消息类型，表示中继是否成功
type RelayResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	NextHop              string   `protobuf:"bytes,2,opt,name=next_hop,json=nextHop,proto3" json:"next_hop,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelayResponse) Reset()         { *m = RelayResponse{} }
func (m *RelayResponse) String() string { return proto.CompactTextString(m) }
func (*RelayResponse) ProtoMessage()    {}
func (*RelayResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8571034d60397816, []int{2}
}

func (m *RelayResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelayResponse.Unmarshal(m, b)
}
func (m *RelayResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelayResponse.Marshal(b, m, deterministic)
}
func (m *RelayResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayResponse.Merge(m, src)
}
func (m *RelayResponse) XXX_Size() int {
	return xxx_messageInfo_RelayResponse.Size(m)
}
func (m *RelayResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RelayResponse proto.InternalMessageInfo

func (m *RelayResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *RelayResponse) GetNextHop() string {
	if m != nil {
		return m.NextHop
	}
	return ""
}

// 事件请求消息
type EventRequest struct {
	EventData            string   `protobuf:"bytes,1,opt,name=event_data,json=eventData,proto3" json:"event_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventRequest) Reset()         { *m = EventRequest{} }
func (m *EventRequest) String() string { return proto.CompactTextString(m) }
func (*EventRequest) ProtoMessage()    {}
func (*EventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8571034d60397816, []int{3}
}

func (m *EventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventRequest.Unmarshal(m, b)
}
func (m *EventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventRequest.Marshal(b, m, deterministic)
}
func (m *EventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventRequest.Merge(m, src)
}
func (m *EventRequest) XXX_Size() int {
	return xxx_messageInfo_EventRequest.Size(m)
}
func (m *EventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EventRequest proto.InternalMessageInfo

func (m *EventRequest) GetEventData() string {
	if m != nil {
		return m.EventData
	}
	return ""
}

// 事件响应消息
type EventResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventResponse) Reset()         { *m = EventResponse{} }
func (m *EventResponse) String() string { return proto.CompactTextString(m) }
func (*EventResponse) ProtoMessage()    {}
func (*EventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8571034d60397816, []int{4}
}

func (m *EventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventResponse.Unmarshal(m, b)
}
func (m *EventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventResponse.Marshal(b, m, deterministic)
}
func (m *EventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventResponse.Merge(m, src)
}
func (m *EventResponse) XXX_Size() int {
	return xxx_messageInfo_EventResponse.Size(m)
}
func (m *EventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EventResponse proto.InternalMessageInfo

func (m *EventResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *EventResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*IDU)(nil), "proto.IDU")
	proto.RegisterType((*RelayRequest)(nil), "proto.RelayRequest")
	proto.RegisterType((*RelayResponse)(nil), "proto.RelayResponse")
	proto.RegisterType((*EventRequest)(nil), "proto.EventRequest")
	proto.RegisterType((*EventResponse)(nil), "proto.EventResponse")
}

func init() { proto.RegisterFile("network.proto", fileDescriptor_8571034d60397816) }

var fileDescriptor_8571034d60397816 = []byte{
	// 383 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0xb1, 0x8e, 0xd3, 0x40,
	0x10, 0x55, 0x08, 0x77, 0x8e, 0xe7, 0x1c, 0x4e, 0x5a, 0xae, 0x58, 0x4e, 0x9c, 0x44, 0xdc, 0x40,
	0x01, 0x41, 0x0a, 0xa2, 0xa3, 0x81, 0x38, 0x88, 0x48, 0x88, 0xc2, 0x22, 0x4d, 0x1a, 0x6b, 0x89,
	0xc7, 0x64, 0x45, 0xbc, 0xb3, 0x78, 0x27, 0x40, 0x3e, 0x97, 0x3f, 0x41, 0x5e, 0xaf, 0x43, 0x52,
	0x5d, 0x65, 0xbf, 0xf7, 0x46, 0xf3, 0xe6, 0xbd, 0x85, 0xb1, 0x41, 0xfe, 0x4d, 0xcd, 0x8f, 0xa9,
	0x6d, 0x88, 0x49, 0x5c, 0xf8, 0x4f, 0xfa, 0x77, 0x00, 0xc3, 0x65, 0xb6, 0x12, 0x13, 0x48, 0xaa,
	0x86, 0xea, 0x42, 0x95, 0x65, 0x83, 0xce, 0xc9, 0xc1, 0xb3, 0xc1, 0x8b, 0x38, 0xbf, 0x6a, 0xb9,
	0xf7, 0x1d, 0x25, 0xee, 0x00, 0x98, 0x8e, 0x03, 0x0f, 0xfc, 0x40, 0xcc, 0xd4, 0xcb, 0x13, 0x48,
	0x36, 0x64, 0x18, 0x0d, 0x17, 0x7c, 0xb0, 0x28, 0x87, 0xdd, 0x86, 0xc0, 0x7d, 0x3d, 0x58, 0x14,
	0x12, 0xa2, 0x00, 0xe5, 0x43, 0xaf, 0xf6, 0x50, 0x3c, 0x85, 0x98, 0x75, 0x8d, 0x8e, 0x55, 0x6d,
	0xe5, 0x45, 0x58, 0xdd, 0x13, 0xe2, 0x06, 0xda, 0x6b, 0xa9, 0x92, 0x97, 0x5e, 0xe9, 0x80, 0x78,
	0x0e, 0xd7, 0xaa, 0x2c, 0x35, 0x6b, 0x32, 0x6a, 0x57, 0x68, 0x53, 0x91, 0x8c, 0xbc, 0xfe, 0xe8,
	0x3f, 0xbd, 0x34, 0x15, 0xa5, 0x2f, 0x21, 0xc9, 0x71, 0xa7, 0x0e, 0x39, 0xfe, 0xdc, 0xa3, 0x6b,
	0xcd, 0x86, 0xba, 0xdc, 0xfb, 0x88, 0x57, 0x33, 0xe8, 0xfa, 0x98, 0x2e, 0xb3, 0x55, 0xde, 0xd2,
	0x69, 0x06, 0xe3, 0x30, 0xed, 0x2c, 0x19, 0xe7, 0xaf, 0x76, 0xfb, 0xcd, 0xa6, 0x6f, 0x65, 0x94,
	0xf7, 0x50, 0x3c, 0x81, 0x91, 0xc1, 0x3f, 0x5c, 0x6c, 0xc9, 0x86, 0x3e, 0xa2, 0x16, 0x7f, 0x22,
	0x9b, 0xbe, 0x82, 0x64, 0xf1, 0x0b, 0x0d, 0xf7, 0x9e, 0x77, 0x00, 0xd8, 0xe2, 0xa2, 0x54, 0xac,
	0x42, 0xbb, 0xb1, 0x67, 0x32, 0xc5, 0x2a, 0x9d, 0xc3, 0x38, 0x8c, 0xdf, 0x6b, 0x2a, 0x21, 0xaa,
	0xd1, 0x39, 0xf5, 0x1d, 0x7b, 0xcf, 0x00, 0x67, 0x8b, 0x90, 0xf3, 0x4b, 0xf7, 0xd0, 0xe2, 0x2d,
	0x8c, 0x3c, 0x6e, 0xdf, 0xf7, 0x71, 0x88, 0x79, 0x5a, 0xc4, 0xed, 0xcd, 0x39, 0xd9, 0x59, 0xcf,
	0xd6, 0x70, 0xbb, 0x30, 0xa5, 0x25, 0x6d, 0x78, 0xbe, 0x55, 0xda, 0xf8, 0xc3, 0x3e, 0x6b, 0xc7,
	0x68, 0xb0, 0x11, 0xef, 0xe0, 0xba, 0xfb, 0xff, 0x48, 0x8d, 0x57, 0xdc, 0x71, 0xf7, 0x69, 0xe0,
	0xe3, 0xee, 0xb3, 0x58, 0x1f, 0x92, 0x35, 0x30, 0x3a, 0x7e, 0xed, 0xb5, 0x6f, 0x97, 0xfe, 0xf3,
	0xe6, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x20, 0xc9, 0x8e, 0xfe, 0x9b, 0x02, 0x00, 0x00,
}
