// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node.proto

package node

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type Peer struct {
	PrivIp               string   `protobuf:"bytes,1,opt,name=priv_ip,json=privIp,proto3" json:"priv_ip,omitempty"`
	PubIp                string   `protobuf:"bytes,2,opt,name=pub_ip,json=pubIp,proto3" json:"pub_ip,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Friend               string   `protobuf:"bytes,4,opt,name=friend,proto3" json:"friend,omitempty"`
	Dialer               bool     `protobuf:"varint,5,opt,name=Dialer,proto3" json:"Dialer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

func (m *Peer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Peer.Unmarshal(m, b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
}
func (m *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(m, src)
}
func (m *Peer) XXX_Size() int {
	return xxx_messageInfo_Peer.Size(m)
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetPrivIp() string {
	if m != nil {
		return m.PrivIp
	}
	return ""
}

func (m *Peer) GetPubIp() string {
	if m != nil {
		return m.PubIp
	}
	return ""
}

func (m *Peer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Peer) GetFriend() string {
	if m != nil {
		return m.Friend
	}
	return ""
}

func (m *Peer) GetDialer() bool {
	if m != nil {
		return m.Dialer
	}
	return false
}

type NodeList struct {
	Nodes                []*Node  `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeList) Reset()         { *m = NodeList{} }
func (m *NodeList) String() string { return proto.CompactTextString(m) }
func (*NodeList) ProtoMessage()    {}
func (*NodeList) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{1}
}

func (m *NodeList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeList.Unmarshal(m, b)
}
func (m *NodeList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeList.Marshal(b, m, deterministic)
}
func (m *NodeList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeList.Merge(m, src)
}
func (m *NodeList) XXX_Size() int {
	return xxx_messageInfo_NodeList.Size(m)
}
func (m *NodeList) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeList.DiscardUnknown(m)
}

var xxx_messageInfo_NodeList proto.InternalMessageInfo

func (m *NodeList) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type Response struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type Request struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{3}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

type Node struct {
	Pub_IP               string   `protobuf:"bytes,1,opt,name=pub_IP,json=pubIP,proto3" json:"pub_IP,omitempty"`
	Grpc_IP              string   `protobuf:"bytes,2,opt,name=grpc_IP,json=grpcIP,proto3" json:"grpc_IP,omitempty"`
	NotifyOthers         int32    `protobuf:"varint,3,opt,name=notify_others,json=notifyOthers,proto3" json:"notify_others,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{4}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetPub_IP() string {
	if m != nil {
		return m.Pub_IP
	}
	return ""
}

func (m *Node) GetGrpc_IP() string {
	if m != nil {
		return m.Grpc_IP
	}
	return ""
}

func (m *Node) GetNotifyOthers() int32 {
	if m != nil {
		return m.NotifyOthers
	}
	return 0
}

func init() {
	proto.RegisterType((*Peer)(nil), "node.Peer")
	proto.RegisterType((*NodeList)(nil), "node.NodeList")
	proto.RegisterType((*Response)(nil), "node.Response")
	proto.RegisterType((*Request)(nil), "node.Request")
	proto.RegisterType((*Node)(nil), "node.Node")
}

func init() { proto.RegisterFile("node.proto", fileDescriptor_0c843d59d2d938e7) }

var fileDescriptor_0c843d59d2d938e7 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xcd, 0xce, 0xda, 0x30,
	0x10, 0x24, 0x25, 0x09, 0xb0, 0x40, 0x0f, 0x96, 0xda, 0x46, 0x9c, 0x22, 0xf7, 0x82, 0x54, 0x84,
	0x54, 0xfa, 0x04, 0xfd, 0x39, 0x34, 0x12, 0x6a, 0x23, 0xf7, 0xd8, 0x03, 0x32, 0x64, 0x01, 0xab,
	0x21, 0x76, 0x6d, 0x27, 0x52, 0x79, 0xe8, 0x3e, 0x43, 0x65, 0x87, 0xa8, 0x28, 0xfa, 0x90, 0xbe,
	0xdb, 0xce, 0x78, 0xd6, 0x3b, 0xb3, 0x5a, 0x80, 0x4a, 0x16, 0xb8, 0x56, 0x5a, 0x5a, 0x49, 0x42,
	0x57, 0xd3, 0x2b, 0x84, 0x39, 0xa2, 0x26, 0x6f, 0x60, 0xa4, 0xb4, 0x68, 0x76, 0x42, 0x25, 0x41,
	0x1a, 0x2c, 0x27, 0x2c, 0x76, 0x30, 0x53, 0xe4, 0x15, 0xc4, 0xaa, 0xde, 0x3b, 0xfe, 0x85, 0xe7,
	0x23, 0x55, 0xef, 0x33, 0x45, 0x08, 0x84, 0x15, 0xbf, 0x60, 0x32, 0xf4, 0xa4, 0xaf, 0xc9, 0x6b,
	0x88, 0x8f, 0x5a, 0x60, 0x55, 0x24, 0x61, 0xfb, 0x45, 0x8b, 0x1c, 0xff, 0x45, 0xf0, 0x12, 0x75,
	0x12, 0xa5, 0xc1, 0x72, 0xcc, 0x6e, 0x88, 0xae, 0x60, 0xfc, 0x4d, 0x16, 0xb8, 0x15, 0xc6, 0x92,
	0x14, 0x22, 0xe7, 0xc7, 0x24, 0x41, 0x3a, 0x5c, 0x4e, 0x37, 0xb0, 0xf6, 0x4e, 0xdd, 0x33, 0x6b,
	0x1f, 0x28, 0x85, 0x31, 0x43, 0xa3, 0x64, 0x65, 0xfc, 0x24, 0x63, 0xb9, 0xad, 0x4d, 0x67, 0xb6,
	0x45, 0x74, 0x02, 0x23, 0x86, 0xbf, 0x6b, 0x34, 0x96, 0xfe, 0x84, 0xd0, 0x75, 0x77, 0xfe, 0xb3,
	0xfc, 0x26, 0xf5, 0xfe, 0x73, 0x97, 0xf7, 0xa4, 0xd5, 0xc1, 0xf1, 0x6d, 0xae, 0xd8, 0xc1, 0x2c,
	0x27, 0x6f, 0x61, 0x5e, 0x49, 0x2b, 0x8e, 0x7f, 0x76, 0xd2, 0x9e, 0x51, 0x1b, 0x9f, 0x30, 0x62,
	0xb3, 0x96, 0xfc, 0xee, 0xb9, 0x8d, 0x84, 0xd9, 0x56, 0xf2, 0xe2, 0x13, 0x2f, 0x79, 0x75, 0x40,
	0x4d, 0x56, 0x30, 0x63, 0x78, 0x12, 0xc6, 0xa2, 0xf6, 0x43, 0xef, 0xec, 0x2f, 0x5e, 0xfe, 0xaf,
	0x5d, 0x52, 0x3a, 0x20, 0xef, 0x81, 0x30, 0xac, 0x0a, 0x6c, 0x64, 0x6d, 0xae, 0x3f, 0x50, 0x37,
	0xa8, 0xb3, 0x9c, 0xcc, 0x5b, 0xdd, 0xcd, 0xff, 0xe2, 0xee, 0x0b, 0x3a, 0xd8, 0xfc, 0x0d, 0x00,
	0x7c, 0xcf, 0xd5, 0x35, 0x91, 0x77, 0x30, 0xf5, 0xab, 0xc1, 0x8b, 0x6c, 0x78, 0xf9, 0xd4, 0xb8,
	0x6e, 0x55, 0x74, 0xd0, 0x89, 0x3f, 0x6a, 0x2d, 0x9e, 0x25, 0x76, 0xf7, 0xd0, 0x13, 0x3b, 0xea,
	0xb1, 0xb8, 0x67, 0xe3, 0x81, 0x78, 0x05, 0xd3, 0xaf, 0xc8, 0x4b, 0x7b, 0xfe, 0x7c, 0xc6, 0xc3,
	0xaf, 0x7e, 0xdc, 0x9e, 0x7e, 0x1f, 0xfb, 0x23, 0xfd, 0xf0, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x14,
	0xcb, 0xaa, 0x70, 0xb2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LoadBalancerClient is the client API for LoadBalancer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LoadBalancerClient interface {
	RegisterNode(ctx context.Context, in *Node, opts ...grpc.CallOption) (*NodeList, error)
	RendevouszServerIP(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Node, error)
}

type loadBalancerClient struct {
	cc *grpc.ClientConn
}

func NewLoadBalancerClient(cc *grpc.ClientConn) LoadBalancerClient {
	return &loadBalancerClient{cc}
}

func (c *loadBalancerClient) RegisterNode(ctx context.Context, in *Node, opts ...grpc.CallOption) (*NodeList, error) {
	out := new(NodeList)
	err := c.cc.Invoke(ctx, "/node.LoadBalancer/RegisterNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loadBalancerClient) RendevouszServerIP(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Node, error) {
	out := new(Node)
	err := c.cc.Invoke(ctx, "/node.LoadBalancer/RendevouszServerIP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoadBalancerServer is the server API for LoadBalancer service.
type LoadBalancerServer interface {
	RegisterNode(context.Context, *Node) (*NodeList, error)
	RendevouszServerIP(context.Context, *Request) (*Node, error)
}

func RegisterLoadBalancerServer(s *grpc.Server, srv LoadBalancerServer) {
	s.RegisterService(&_LoadBalancer_serviceDesc, srv)
}

func _LoadBalancer_RegisterNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadBalancerServer).RegisterNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.LoadBalancer/RegisterNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadBalancerServer).RegisterNode(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoadBalancer_RendevouszServerIP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadBalancerServer).RendevouszServerIP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.LoadBalancer/RendevouszServerIP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadBalancerServer).RendevouszServerIP(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _LoadBalancer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "node.LoadBalancer",
	HandlerType: (*LoadBalancerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNode",
			Handler:    _LoadBalancer_RegisterNode_Handler,
		},
		{
			MethodName: "RendevouszServerIP",
			Handler:    _LoadBalancer_RendevouszServerIP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}

// RendezvousClient is the client API for Rendezvous service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RendezvousClient interface {
	NodeRemoval(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Response, error)
	NodeArrival(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Response, error)
	PeerArrival(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*Response, error)
	PeerRemoval(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*Response, error)
	HealthCheck(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type rendezvousClient struct {
	cc *grpc.ClientConn
}

func NewRendezvousClient(cc *grpc.ClientConn) RendezvousClient {
	return &rendezvousClient{cc}
}

func (c *rendezvousClient) NodeRemoval(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/node.Rendezvous/NodeRemoval", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendezvousClient) NodeArrival(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/node.Rendezvous/NodeArrival", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendezvousClient) PeerArrival(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/node.Rendezvous/PeerArrival", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendezvousClient) PeerRemoval(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/node.Rendezvous/PeerRemoval", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rendezvousClient) HealthCheck(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/node.Rendezvous/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RendezvousServer is the server API for Rendezvous service.
type RendezvousServer interface {
	NodeRemoval(context.Context, *Node) (*Response, error)
	NodeArrival(context.Context, *Node) (*Response, error)
	PeerArrival(context.Context, *Peer) (*Response, error)
	PeerRemoval(context.Context, *Peer) (*Response, error)
	HealthCheck(context.Context, *Request) (*Response, error)
}

func RegisterRendezvousServer(s *grpc.Server, srv RendezvousServer) {
	s.RegisterService(&_Rendezvous_serviceDesc, srv)
}

func _Rendezvous_NodeRemoval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendezvousServer).NodeRemoval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.Rendezvous/NodeRemoval",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendezvousServer).NodeRemoval(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rendezvous_NodeArrival_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendezvousServer).NodeArrival(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.Rendezvous/NodeArrival",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendezvousServer).NodeArrival(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rendezvous_PeerArrival_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendezvousServer).PeerArrival(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.Rendezvous/PeerArrival",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendezvousServer).PeerArrival(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rendezvous_PeerRemoval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendezvousServer).PeerRemoval(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.Rendezvous/PeerRemoval",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendezvousServer).PeerRemoval(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rendezvous_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RendezvousServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.Rendezvous/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RendezvousServer).HealthCheck(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Rendezvous_serviceDesc = grpc.ServiceDesc{
	ServiceName: "node.Rendezvous",
	HandlerType: (*RendezvousServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NodeRemoval",
			Handler:    _Rendezvous_NodeRemoval_Handler,
		},
		{
			MethodName: "NodeArrival",
			Handler:    _Rendezvous_NodeArrival_Handler,
		},
		{
			MethodName: "PeerArrival",
			Handler:    _Rendezvous_PeerArrival_Handler,
		},
		{
			MethodName: "PeerRemoval",
			Handler:    _Rendezvous_PeerRemoval_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _Rendezvous_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}
