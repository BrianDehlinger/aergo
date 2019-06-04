// Code generated by protoc-gen-go. DO NOT EDIT.
// source: raft.proto

package types

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// cluster member for raft consensus
type MembershipChangeType int32

const (
	MembershipChangeType_ADD_MEMBER    MembershipChangeType = 0
	MembershipChangeType_REMOVE_MEMBER MembershipChangeType = 1
)

var MembershipChangeType_name = map[int32]string{
	0: "ADD_MEMBER",
	1: "REMOVE_MEMBER",
}

var MembershipChangeType_value = map[string]int32{
	"ADD_MEMBER":    0,
	"REMOVE_MEMBER": 1,
}

func (x MembershipChangeType) String() string {
	return proto.EnumName(MembershipChangeType_name, int32(x))
}

func (MembershipChangeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{0}
}

type MemberAttr struct {
	ID                   uint64   `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Url                  string   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	PeerID               []byte   `protobuf:"bytes,4,opt,name=peerID,proto3" json:"peerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MemberAttr) Reset()         { *m = MemberAttr{} }
func (m *MemberAttr) String() string { return proto.CompactTextString(m) }
func (*MemberAttr) ProtoMessage()    {}
func (*MemberAttr) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{0}
}

func (m *MemberAttr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MemberAttr.Unmarshal(m, b)
}
func (m *MemberAttr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MemberAttr.Marshal(b, m, deterministic)
}
func (m *MemberAttr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MemberAttr.Merge(m, src)
}
func (m *MemberAttr) XXX_Size() int {
	return xxx_messageInfo_MemberAttr.Size(m)
}
func (m *MemberAttr) XXX_DiscardUnknown() {
	xxx_messageInfo_MemberAttr.DiscardUnknown(m)
}

var xxx_messageInfo_MemberAttr proto.InternalMessageInfo

func (m *MemberAttr) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *MemberAttr) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MemberAttr) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *MemberAttr) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

type MembershipChange struct {
	Type                 MembershipChangeType `protobuf:"varint,1,opt,name=type,proto3,enum=types.MembershipChangeType" json:"type,omitempty"`
	Attr                 *MemberAttr          `protobuf:"bytes,2,opt,name=attr,proto3" json:"attr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *MembershipChange) Reset()         { *m = MembershipChange{} }
func (m *MembershipChange) String() string { return proto.CompactTextString(m) }
func (*MembershipChange) ProtoMessage()    {}
func (*MembershipChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{1}
}

func (m *MembershipChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MembershipChange.Unmarshal(m, b)
}
func (m *MembershipChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MembershipChange.Marshal(b, m, deterministic)
}
func (m *MembershipChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MembershipChange.Merge(m, src)
}
func (m *MembershipChange) XXX_Size() int {
	return xxx_messageInfo_MembershipChange.Size(m)
}
func (m *MembershipChange) XXX_DiscardUnknown() {
	xxx_messageInfo_MembershipChange.DiscardUnknown(m)
}

var xxx_messageInfo_MembershipChange proto.InternalMessageInfo

func (m *MembershipChange) GetType() MembershipChangeType {
	if m != nil {
		return m.Type
	}
	return MembershipChangeType_ADD_MEMBER
}

func (m *MembershipChange) GetAttr() *MemberAttr {
	if m != nil {
		return m.Attr
	}
	return nil
}

type MembershipChangeReply struct {
	Attr                 *MemberAttr `protobuf:"bytes,1,opt,name=attr,proto3" json:"attr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *MembershipChangeReply) Reset()         { *m = MembershipChangeReply{} }
func (m *MembershipChangeReply) String() string { return proto.CompactTextString(m) }
func (*MembershipChangeReply) ProtoMessage()    {}
func (*MembershipChangeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{2}
}

func (m *MembershipChangeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MembershipChangeReply.Unmarshal(m, b)
}
func (m *MembershipChangeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MembershipChangeReply.Marshal(b, m, deterministic)
}
func (m *MembershipChangeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MembershipChangeReply.Merge(m, src)
}
func (m *MembershipChangeReply) XXX_Size() int {
	return xxx_messageInfo_MembershipChangeReply.Size(m)
}
func (m *MembershipChangeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MembershipChangeReply.DiscardUnknown(m)
}

var xxx_messageInfo_MembershipChangeReply proto.InternalMessageInfo

func (m *MembershipChangeReply) GetAttr() *MemberAttr {
	if m != nil {
		return m.Attr
	}
	return nil
}

// data types for raft support
// GetClusterInfoRequest
type GetClusterInfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetClusterInfoRequest) Reset()         { *m = GetClusterInfoRequest{} }
func (m *GetClusterInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetClusterInfoRequest) ProtoMessage()    {}
func (*GetClusterInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{3}
}

func (m *GetClusterInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetClusterInfoRequest.Unmarshal(m, b)
}
func (m *GetClusterInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetClusterInfoRequest.Marshal(b, m, deterministic)
}
func (m *GetClusterInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetClusterInfoRequest.Merge(m, src)
}
func (m *GetClusterInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetClusterInfoRequest.Size(m)
}
func (m *GetClusterInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetClusterInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetClusterInfoRequest proto.InternalMessageInfo

type GetClusterInfoResponse struct {
	ChainID              []byte        `protobuf:"bytes,1,opt,name=chainID,proto3" json:"chainID,omitempty"`
	Error                string        `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	MbrAttrs             []*MemberAttr `protobuf:"bytes,3,rep,name=mbrAttrs,proto3" json:"mbrAttrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetClusterInfoResponse) Reset()         { *m = GetClusterInfoResponse{} }
func (m *GetClusterInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetClusterInfoResponse) ProtoMessage()    {}
func (*GetClusterInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b042552c306ae59b, []int{4}
}

func (m *GetClusterInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetClusterInfoResponse.Unmarshal(m, b)
}
func (m *GetClusterInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetClusterInfoResponse.Marshal(b, m, deterministic)
}
func (m *GetClusterInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetClusterInfoResponse.Merge(m, src)
}
func (m *GetClusterInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetClusterInfoResponse.Size(m)
}
func (m *GetClusterInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetClusterInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetClusterInfoResponse proto.InternalMessageInfo

func (m *GetClusterInfoResponse) GetChainID() []byte {
	if m != nil {
		return m.ChainID
	}
	return nil
}

func (m *GetClusterInfoResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *GetClusterInfoResponse) GetMbrAttrs() []*MemberAttr {
	if m != nil {
		return m.MbrAttrs
	}
	return nil
}

func init() {
	proto.RegisterEnum("types.MembershipChangeType", MembershipChangeType_name, MembershipChangeType_value)
	proto.RegisterType((*MemberAttr)(nil), "types.MemberAttr")
	proto.RegisterType((*MembershipChange)(nil), "types.MembershipChange")
	proto.RegisterType((*MembershipChangeReply)(nil), "types.MembershipChangeReply")
	proto.RegisterType((*GetClusterInfoRequest)(nil), "types.GetClusterInfoRequest")
	proto.RegisterType((*GetClusterInfoResponse)(nil), "types.GetClusterInfoResponse")
}

func init() { proto.RegisterFile("raft.proto", fileDescriptor_b042552c306ae59b) }

var fileDescriptor_b042552c306ae59b = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0x4d, 0x6b, 0xc2, 0x40,
	0x14, 0x6c, 0x4c, 0xb4, 0xf5, 0xd5, 0x8a, 0x2e, 0x6a, 0x43, 0x0b, 0x25, 0x04, 0x0a, 0x52, 0x68,
	0x04, 0x7b, 0xea, 0xa5, 0xa0, 0x26, 0x94, 0x1c, 0x42, 0x61, 0x29, 0x3d, 0x78, 0x29, 0x89, 0x3c,
	0x4d, 0x8a, 0xc9, 0x6e, 0x77, 0x37, 0x14, 0xff, 0x7d, 0x71, 0xa3, 0x42, 0x45, 0x7a, 0xda, 0xf7,
	0x31, 0xf3, 0x66, 0x98, 0x05, 0x10, 0xf1, 0x52, 0x79, 0x5c, 0x30, 0xc5, 0x48, 0x5d, 0x6d, 0x38,
	0xca, 0x9b, 0x26, 0x1f, 0xf3, 0x6a, 0xe2, 0xce, 0x01, 0x22, 0xcc, 0x13, 0x14, 0x13, 0xa5, 0x04,
	0x69, 0x43, 0x2d, 0xf4, 0x6d, 0xc3, 0x31, 0x86, 0x16, 0xad, 0x85, 0x3e, 0x21, 0x60, 0x15, 0x71,
	0x8e, 0x76, 0xcd, 0x31, 0x86, 0x4d, 0xaa, 0x6b, 0xd2, 0x01, 0xb3, 0x14, 0x6b, 0xdb, 0xd4, 0xa3,
	0x6d, 0x49, 0x06, 0xd0, 0xe0, 0x88, 0x22, 0xf4, 0x6d, 0xcb, 0x31, 0x86, 0x2d, 0xba, 0xeb, 0xdc,
	0x2f, 0xe8, 0x54, 0xb7, 0x65, 0x9a, 0xf1, 0x59, 0x1a, 0x17, 0x2b, 0x24, 0x23, 0xb0, 0xb6, 0x1e,
	0xb4, 0x46, 0x7b, 0x7c, 0xeb, 0x69, 0x43, 0xde, 0x31, 0xec, 0x7d, 0xc3, 0x91, 0x6a, 0x20, 0xb9,
	0x07, 0x2b, 0x56, 0x4a, 0x68, 0x0b, 0x97, 0xe3, 0xee, 0x1f, 0xc2, 0xd6, 0x33, 0xd5, 0x6b, 0xf7,
	0x05, 0xfa, 0xc7, 0x47, 0x28, 0xf2, 0xf5, 0xe6, 0xc0, 0x37, 0xfe, 0xe7, 0x5f, 0x43, 0xff, 0x15,
	0xd5, 0x6c, 0x5d, 0x4a, 0x85, 0x22, 0x2c, 0x96, 0x8c, 0xe2, 0x77, 0x89, 0x52, 0xb9, 0x3f, 0x30,
	0x38, 0x5e, 0x48, 0xce, 0x0a, 0x89, 0xc4, 0x86, 0xf3, 0x45, 0x1a, 0x67, 0xc5, 0x2e, 0xb1, 0x16,
	0xdd, 0xb7, 0xa4, 0x07, 0x75, 0x14, 0x82, 0x89, 0x5d, 0x6e, 0x55, 0x43, 0x1e, 0xe1, 0x22, 0x4f,
	0xb4, 0xa6, 0xb4, 0x4d, 0xc7, 0x3c, 0xed, 0xe6, 0x00, 0x79, 0x78, 0x86, 0xde, 0xa9, 0x58, 0x48,
	0x1b, 0x60, 0xe2, 0xfb, 0x9f, 0x51, 0x10, 0x4d, 0x03, 0xda, 0x39, 0x23, 0x5d, 0xb8, 0xa2, 0x41,
	0xf4, 0xf6, 0x11, 0xec, 0x47, 0xc6, 0xd4, 0x99, 0xdf, 0xad, 0x32, 0x95, 0x96, 0x89, 0xb7, 0x60,
	0xf9, 0x28, 0x46, 0xb1, 0x62, 0x19, 0xab, 0xde, 0x91, 0x56, 0x4c, 0x1a, 0xfa, 0xf7, 0x9f, 0x7e,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x01, 0xb2, 0x9b, 0x00, 0x1d, 0x02, 0x00, 0x00,
}
