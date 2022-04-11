// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: kvm.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// KVMClient is the client API for KVM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVMClient interface {
	RegisterAgentIdentity(ctx context.Context, in *AgentIdentity, opts ...grpc.CallOption) (*Status, error)
	SendPolicy(ctx context.Context, opts ...grpc.CallOption) (KVM_SendPolicyClient, error)
}

type kVMClient struct {
	cc grpc.ClientConnInterface
}

func NewKVMClient(cc grpc.ClientConnInterface) KVMClient {
	return &kVMClient{cc}
}

func (c *kVMClient) RegisterAgentIdentity(ctx context.Context, in *AgentIdentity, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/kvm.KVM/registerAgentIdentity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVMClient) SendPolicy(ctx context.Context, opts ...grpc.CallOption) (KVM_SendPolicyClient, error) {
	stream, err := c.cc.NewStream(ctx, &KVM_ServiceDesc.Streams[0], "/kvm.KVM/sendPolicy", opts...)
	if err != nil {
		return nil, err
	}
	x := &kVMSendPolicyClient{stream}
	return x, nil
}

type KVM_SendPolicyClient interface {
	Send(*Status) error
	Recv() (*PolicyData, error)
	grpc.ClientStream
}

type kVMSendPolicyClient struct {
	grpc.ClientStream
}

func (x *kVMSendPolicyClient) Send(m *Status) error {
	return x.ClientStream.SendMsg(m)
}

func (x *kVMSendPolicyClient) Recv() (*PolicyData, error) {
	m := new(PolicyData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// KVMServer is the server API for KVM service.
// All implementations should embed UnimplementedKVMServer
// for forward compatibility
type KVMServer interface {
	RegisterAgentIdentity(context.Context, *AgentIdentity) (*Status, error)
	SendPolicy(KVM_SendPolicyServer) error
}

// UnimplementedKVMServer should be embedded to have forward compatible implementations.
type UnimplementedKVMServer struct {
}

func (UnimplementedKVMServer) RegisterAgentIdentity(context.Context, *AgentIdentity) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAgentIdentity not implemented")
}
func (UnimplementedKVMServer) SendPolicy(KVM_SendPolicyServer) error {
	return status.Errorf(codes.Unimplemented, "method SendPolicy not implemented")
}

// UnsafeKVMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVMServer will
// result in compilation errors.
type UnsafeKVMServer interface {
	mustEmbedUnimplementedKVMServer()
}

func RegisterKVMServer(s grpc.ServiceRegistrar, srv KVMServer) {
	s.RegisterService(&KVM_ServiceDesc, srv)
}

func _KVM_RegisterAgentIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVMServer).RegisterAgentIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvm.KVM/registerAgentIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVMServer).RegisterAgentIdentity(ctx, req.(*AgentIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVM_SendPolicy_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(KVMServer).SendPolicy(&kVMSendPolicyServer{stream})
}

type KVM_SendPolicyServer interface {
	Send(*PolicyData) error
	Recv() (*Status, error)
	grpc.ServerStream
}

type kVMSendPolicyServer struct {
	grpc.ServerStream
}

func (x *kVMSendPolicyServer) Send(m *PolicyData) error {
	return x.ServerStream.SendMsg(m)
}

func (x *kVMSendPolicyServer) Recv() (*Status, error) {
	m := new(Status)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// KVM_ServiceDesc is the grpc.ServiceDesc for KVM service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KVM_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kvm.KVM",
	HandlerType: (*KVMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "registerAgentIdentity",
			Handler:    _KVM_RegisterAgentIdentity_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "sendPolicy",
			Handler:       _KVM_SendPolicy_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "kvm.proto",
}
