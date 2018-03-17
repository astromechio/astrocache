// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	ReserveIDRequest
	ReserveIDResponse
	NewNodeRequest
	NewNodeResponse
	GetChainRequest
	GetChainResponse
	AddBlockRequest
	AddBlockResponse
	Node
*/
package service

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type ReserveIDRequest struct {
	ProposingNID string `protobuf:"bytes,1,opt,name=ProposingNID" json:"ProposingNID,omitempty"`
}

func (m *ReserveIDRequest) Reset()                    { *m = ReserveIDRequest{} }
func (m *ReserveIDRequest) String() string            { return proto.CompactTextString(m) }
func (*ReserveIDRequest) ProtoMessage()               {}
func (*ReserveIDRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ReserveIDRequest) GetProposingNID() string {
	if m != nil {
		return m.ProposingNID
	}
	return ""
}

type ReserveIDResponse struct {
	BlockID string `protobuf:"bytes,1,opt,name=BlockID" json:"BlockID,omitempty"`
}

func (m *ReserveIDResponse) Reset()                    { *m = ReserveIDResponse{} }
func (m *ReserveIDResponse) String() string            { return proto.CompactTextString(m) }
func (*ReserveIDResponse) ProtoMessage()               {}
func (*ReserveIDResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ReserveIDResponse) GetBlockID() string {
	if m != nil {
		return m.BlockID
	}
	return ""
}

type NewNodeRequest struct {
	Node     *Node  `protobuf:"bytes,1,opt,name=Node" json:"Node,omitempty"`
	JoinCode string `protobuf:"bytes,2,opt,name=JoinCode" json:"JoinCode,omitempty"`
}

func (m *NewNodeRequest) Reset()                    { *m = NewNodeRequest{} }
func (m *NewNodeRequest) String() string            { return proto.CompactTextString(m) }
func (*NewNodeRequest) ProtoMessage()               {}
func (*NewNodeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *NewNodeRequest) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *NewNodeRequest) GetJoinCode() string {
	if m != nil {
		return m.JoinCode
	}
	return ""
}

type NewNodeResponse struct {
	EncGlobalKey []byte `protobuf:"bytes,1,opt,name=EncGlobalKey,proto3" json:"EncGlobalKey,omitempty"`
	Master       *Node  `protobuf:"bytes,2,opt,name=Master" json:"Master,omitempty"`
	Verifier     *Node  `protobuf:"bytes,3,opt,name=Verifier" json:"Verifier,omitempty"`
	IsPrimary    bool   `protobuf:"varint,4,opt,name=IsPrimary" json:"IsPrimary,omitempty"`
}

func (m *NewNodeResponse) Reset()                    { *m = NewNodeResponse{} }
func (m *NewNodeResponse) String() string            { return proto.CompactTextString(m) }
func (*NewNodeResponse) ProtoMessage()               {}
func (*NewNodeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NewNodeResponse) GetEncGlobalKey() []byte {
	if m != nil {
		return m.EncGlobalKey
	}
	return nil
}

func (m *NewNodeResponse) GetMaster() *Node {
	if m != nil {
		return m.Master
	}
	return nil
}

func (m *NewNodeResponse) GetVerifier() *Node {
	if m != nil {
		return m.Verifier
	}
	return nil
}

func (m *NewNodeResponse) GetIsPrimary() bool {
	if m != nil {
		return m.IsPrimary
	}
	return false
}

type GetChainRequest struct {
	After string `protobuf:"bytes,1,opt,name=After" json:"After,omitempty"`
}

func (m *GetChainRequest) Reset()                    { *m = GetChainRequest{} }
func (m *GetChainRequest) String() string            { return proto.CompactTextString(m) }
func (*GetChainRequest) ProtoMessage()               {}
func (*GetChainRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetChainRequest) GetAfter() string {
	if m != nil {
		return m.After
	}
	return ""
}

type GetChainResponse struct {
	Blocks []byte `protobuf:"bytes,1,opt,name=Blocks,proto3" json:"Blocks,omitempty"`
}

func (m *GetChainResponse) Reset()                    { *m = GetChainResponse{} }
func (m *GetChainResponse) String() string            { return proto.CompactTextString(m) }
func (*GetChainResponse) ProtoMessage()               {}
func (*GetChainResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetChainResponse) GetBlocks() []byte {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type AddBlockRequest struct {
	Block        []byte `protobuf:"bytes,1,opt,name=Block,proto3" json:"Block,omitempty"`
	ProposingNID string `protobuf:"bytes,2,opt,name=ProposingNID" json:"ProposingNID,omitempty"`
}

func (m *AddBlockRequest) Reset()                    { *m = AddBlockRequest{} }
func (m *AddBlockRequest) String() string            { return proto.CompactTextString(m) }
func (*AddBlockRequest) ProtoMessage()               {}
func (*AddBlockRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AddBlockRequest) GetBlock() []byte {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *AddBlockRequest) GetProposingNID() string {
	if m != nil {
		return m.ProposingNID
	}
	return ""
}

type AddBlockResponse struct {
}

func (m *AddBlockResponse) Reset()                    { *m = AddBlockResponse{} }
func (m *AddBlockResponse) String() string            { return proto.CompactTextString(m) }
func (*AddBlockResponse) ProtoMessage()               {}
func (*AddBlockResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

// Node defines a node in the network
type Node struct {
	NID       string `protobuf:"bytes,1,opt,name=NID" json:"NID,omitempty"`
	Address   string `protobuf:"bytes,2,opt,name=Address" json:"Address,omitempty"`
	Type      string `protobuf:"bytes,3,opt,name=Type" json:"Type,omitempty"`
	PubKey    []byte `protobuf:"bytes,4,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	ParentNID string `protobuf:"bytes,5,opt,name=ParentNID" json:"ParentNID,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Node) GetNID() string {
	if m != nil {
		return m.NID
	}
	return ""
}

func (m *Node) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Node) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Node) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Node) GetParentNID() string {
	if m != nil {
		return m.ParentNID
	}
	return ""
}

func init() {
	proto.RegisterType((*ReserveIDRequest)(nil), "ReserveIDRequest")
	proto.RegisterType((*ReserveIDResponse)(nil), "ReserveIDResponse")
	proto.RegisterType((*NewNodeRequest)(nil), "NewNodeRequest")
	proto.RegisterType((*NewNodeResponse)(nil), "NewNodeResponse")
	proto.RegisterType((*GetChainRequest)(nil), "GetChainRequest")
	proto.RegisterType((*GetChainResponse)(nil), "GetChainResponse")
	proto.RegisterType((*AddBlockRequest)(nil), "AddBlockRequest")
	proto.RegisterType((*AddBlockResponse)(nil), "AddBlockResponse")
	proto.RegisterType((*Node)(nil), "Node")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Master service

type MasterClient interface {
	RequestReservedID(ctx context.Context, in *ReserveIDRequest, opts ...grpc.CallOption) (*ReserveIDResponse, error)
	AddNode(ctx context.Context, in *NewNodeRequest, opts ...grpc.CallOption) (*NewNodeResponse, error)
	GetChain(ctx context.Context, in *GetChainRequest, opts ...grpc.CallOption) (*GetChainResponse, error)
	AddBlock(ctx context.Context, in *AddBlockRequest, opts ...grpc.CallOption) (*AddBlockResponse, error)
}

type masterClient struct {
	cc *grpc.ClientConn
}

func NewMasterClient(cc *grpc.ClientConn) MasterClient {
	return &masterClient{cc}
}

func (c *masterClient) RequestReservedID(ctx context.Context, in *ReserveIDRequest, opts ...grpc.CallOption) (*ReserveIDResponse, error) {
	out := new(ReserveIDResponse)
	err := grpc.Invoke(ctx, "/Master/RequestReservedID", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterClient) AddNode(ctx context.Context, in *NewNodeRequest, opts ...grpc.CallOption) (*NewNodeResponse, error) {
	out := new(NewNodeResponse)
	err := grpc.Invoke(ctx, "/Master/AddNode", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterClient) GetChain(ctx context.Context, in *GetChainRequest, opts ...grpc.CallOption) (*GetChainResponse, error) {
	out := new(GetChainResponse)
	err := grpc.Invoke(ctx, "/Master/GetChain", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterClient) AddBlock(ctx context.Context, in *AddBlockRequest, opts ...grpc.CallOption) (*AddBlockResponse, error) {
	out := new(AddBlockResponse)
	err := grpc.Invoke(ctx, "/Master/AddBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Master service

type MasterServer interface {
	RequestReservedID(context.Context, *ReserveIDRequest) (*ReserveIDResponse, error)
	AddNode(context.Context, *NewNodeRequest) (*NewNodeResponse, error)
	GetChain(context.Context, *GetChainRequest) (*GetChainResponse, error)
	AddBlock(context.Context, *AddBlockRequest) (*AddBlockResponse, error)
}

func RegisterMasterServer(s *grpc.Server, srv MasterServer) {
	s.RegisterService(&_Master_serviceDesc, srv)
}

func _Master_RequestReservedID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReserveIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).RequestReservedID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Master/RequestReservedID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).RequestReservedID(ctx, req.(*ReserveIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Master_AddNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).AddNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Master/AddNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).AddNode(ctx, req.(*NewNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Master_GetChain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).GetChain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Master/GetChain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).GetChain(ctx, req.(*GetChainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Master_AddBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).AddBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Master/AddBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).AddBlock(ctx, req.(*AddBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Master_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Master",
	HandlerType: (*MasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestReservedID",
			Handler:    _Master_RequestReservedID_Handler,
		},
		{
			MethodName: "AddNode",
			Handler:    _Master_AddNode_Handler,
		},
		{
			MethodName: "GetChain",
			Handler:    _Master_GetChain_Handler,
		},
		{
			MethodName: "AddBlock",
			Handler:    _Master_AddBlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 434 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0xd1, 0x8a, 0xd3, 0x40,
	0x14, 0x25, 0xbb, 0x6d, 0x4d, 0xaf, 0xd5, 0x24, 0x83, 0x48, 0x0c, 0x0a, 0xeb, 0xbc, 0xb8, 0x88,
	0xce, 0xc2, 0x0a, 0x3e, 0xf8, 0x56, 0xb7, 0x52, 0xe2, 0x62, 0x09, 0x41, 0x7c, 0x4f, 0x9b, 0xbb,
	0x1a, 0xac, 0x99, 0x38, 0x93, 0x55, 0x0a, 0xfe, 0x86, 0x7f, 0xe6, 0x07, 0xc9, 0xdc, 0x4c, 0x92,
	0x36, 0xdd, 0xb7, 0x9c, 0x33, 0x77, 0xce, 0x9d, 0x7b, 0xce, 0x0d, 0x3c, 0xd0, 0xa8, 0x7e, 0x15,
	0x1b, 0x14, 0x95, 0x92, 0xb5, 0xe4, 0x6f, 0xc1, 0x4f, 0xd1, 0x50, 0x18, 0x2f, 0x52, 0xfc, 0x79,
	0x8b, 0xba, 0x66, 0x1c, 0x66, 0x89, 0x92, 0x95, 0xd4, 0x45, 0xf9, 0x75, 0x15, 0x2f, 0x42, 0xe7,
	0xcc, 0x39, 0x9f, 0xa6, 0x07, 0x1c, 0x7f, 0x0d, 0xc1, 0xde, 0x3d, 0x5d, 0xc9, 0x52, 0x23, 0x0b,
	0xe1, 0xde, 0xfb, 0xad, 0xdc, 0x7c, 0xef, 0xee, 0xb4, 0x90, 0x2f, 0xe1, 0xe1, 0x0a, 0x7f, 0xaf,
	0x64, 0x8e, 0x6d, 0x93, 0x27, 0x30, 0x32, 0x90, 0x0a, 0xef, 0x5f, 0x8e, 0x05, 0x9d, 0x11, 0xc5,
	0x22, 0x70, 0x3f, 0xca, 0xa2, 0xbc, 0x32, 0xc7, 0x27, 0xa4, 0xd3, 0x61, 0xfe, 0xd7, 0x01, 0xaf,
	0x53, 0xb2, 0x6d, 0x39, 0xcc, 0x3e, 0x94, 0x9b, 0xe5, 0x56, 0xae, 0xb3, 0xed, 0x35, 0xee, 0x48,
	0x72, 0x96, 0x1e, 0x70, 0xec, 0x19, 0x4c, 0x3e, 0x65, 0xba, 0x46, 0x45, 0x8a, 0x5d, 0x43, 0x4b,
	0xb2, 0xe7, 0xe0, 0x7e, 0x41, 0x55, 0xdc, 0x14, 0xa8, 0xc2, 0xd3, 0xfd, 0x82, 0x8e, 0x66, 0x4f,
	0x61, 0x1a, 0xeb, 0x44, 0x15, 0x3f, 0x32, 0xb5, 0x0b, 0x47, 0x67, 0xce, 0xb9, 0x9b, 0xf6, 0x04,
	0x7f, 0x01, 0xde, 0x12, 0xeb, 0xab, 0x6f, 0x59, 0x51, 0xb6, 0x13, 0x3e, 0x82, 0xf1, 0xfc, 0xc6,
	0x74, 0x6c, 0xbc, 0x68, 0x00, 0x7f, 0x09, 0x7e, 0x5f, 0x68, 0x07, 0x78, 0x0c, 0x13, 0x32, 0x4a,
	0xdb, 0xa7, 0x5b, 0xc4, 0xaf, 0xc1, 0x9b, 0xe7, 0x39, 0x81, 0x3d, 0x51, 0xc2, 0xb6, 0xb2, 0x01,
	0x47, 0x89, 0x9d, 0xdc, 0x91, 0x18, 0x03, 0xbf, 0x17, 0x6b, 0x1a, 0xf3, 0x3f, 0x4d, 0x08, 0xcc,
	0x87, 0xd3, 0x3e, 0x68, 0xf3, 0x69, 0xa2, 0x9c, 0xe7, 0xb9, 0x42, 0xad, 0xad, 0x58, 0x0b, 0x19,
	0x83, 0xd1, 0xe7, 0x5d, 0x85, 0x64, 0xd3, 0x34, 0xa5, 0x6f, 0x33, 0x40, 0x72, 0xbb, 0x36, 0xde,
	0x8f, 0x9a, 0x01, 0x1a, 0x64, 0x3c, 0x4b, 0x32, 0x85, 0x65, 0x6d, 0xd4, 0xc7, 0x74, 0xa1, 0x27,
	0x2e, 0xff, 0x39, 0x6d, 0x28, 0xec, 0x9d, 0x59, 0x27, 0x9a, 0xd0, 0x6e, 0x55, 0x1e, 0x2f, 0x58,
	0x20, 0x86, 0xab, 0x19, 0x31, 0x71, 0xbc, 0x75, 0xaf, 0xe8, 0xa9, 0x34, 0x87, 0x27, 0x0e, 0xb7,
	0x2c, 0xf2, 0xc5, 0x70, 0x59, 0x2e, 0xc0, 0x6d, 0xfd, 0x67, 0xbe, 0x18, 0x64, 0x16, 0x05, 0xe2,
	0x28, 0x9c, 0x0b, 0x70, 0x5b, 0xdf, 0x98, 0x2f, 0x06, 0x79, 0x44, 0x81, 0x18, 0x9a, 0xba, 0x9e,
	0xd0, 0x9f, 0xf5, 0xe6, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7d, 0xd1, 0x68, 0x5b, 0x6a, 0x03,
	0x00, 0x00,
}