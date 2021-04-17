// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package nodeattestorv1

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

// NodeAttestorClient is the client API for NodeAttestor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeAttestorClient interface {
	// Attest attests attestation payload received from the agent and
	// optionally participates in challenge/response attestation mechanics.
	//
	// The attestation flow is as follows:
	// 1. SPIRE Server opens up a stream to the plugin via Attest.
	// 2. SPIRE Server sends a request containing the attestation payload
	//    received from the agent.
	// 3. Optionally, the plugin responds with a challenge:
	//    3a. SPIRE Server sends the challenge to the agent.
	//    3b. SPIRE Agent responds with the challenge response.
	//    3c. SPIRE Server sends the challenge response to the plugin.
	//    3d. Step 3 is repeated until the plugin is satisfied and does
	//        not respond with an additional challenge.
	// 4. The plugin returns the attestation results to SPIRE Server and closes
	//    the stream.
	Attest(ctx context.Context, opts ...grpc.CallOption) (NodeAttestor_AttestClient, error)
}

type nodeAttestorClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeAttestorClient(cc grpc.ClientConnInterface) NodeAttestorClient {
	return &nodeAttestorClient{cc}
}

func (c *nodeAttestorClient) Attest(ctx context.Context, opts ...grpc.CallOption) (NodeAttestor_AttestClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeAttestor_ServiceDesc.Streams[0], "/spire.plugin.server.nodeattestor.v1.NodeAttestor/Attest", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeAttestorAttestClient{stream}
	return x, nil
}

type NodeAttestor_AttestClient interface {
	Send(*AttestRequest) error
	Recv() (*AttestResponse, error)
	grpc.ClientStream
}

type nodeAttestorAttestClient struct {
	grpc.ClientStream
}

func (x *nodeAttestorAttestClient) Send(m *AttestRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nodeAttestorAttestClient) Recv() (*AttestResponse, error) {
	m := new(AttestResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NodeAttestorServer is the server API for NodeAttestor service.
// All implementations must embed UnimplementedNodeAttestorServer
// for forward compatibility
type NodeAttestorServer interface {
	// Attest attests attestation payload received from the agent and
	// optionally participates in challenge/response attestation mechanics.
	//
	// The attestation flow is as follows:
	// 1. SPIRE Server opens up a stream to the plugin via Attest.
	// 2. SPIRE Server sends a request containing the attestation payload
	//    received from the agent.
	// 3. Optionally, the plugin responds with a challenge:
	//    3a. SPIRE Server sends the challenge to the agent.
	//    3b. SPIRE Agent responds with the challenge response.
	//    3c. SPIRE Server sends the challenge response to the plugin.
	//    3d. Step 3 is repeated until the plugin is satisfied and does
	//        not respond with an additional challenge.
	// 4. The plugin returns the attestation results to SPIRE Server and closes
	//    the stream.
	Attest(NodeAttestor_AttestServer) error
	mustEmbedUnimplementedNodeAttestorServer()
}

// UnimplementedNodeAttestorServer must be embedded to have forward compatible implementations.
type UnimplementedNodeAttestorServer struct {
}

func (UnimplementedNodeAttestorServer) Attest(NodeAttestor_AttestServer) error {
	return status.Errorf(codes.Unimplemented, "method Attest not implemented")
}
func (UnimplementedNodeAttestorServer) mustEmbedUnimplementedNodeAttestorServer() {}

// UnsafeNodeAttestorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeAttestorServer will
// result in compilation errors.
type UnsafeNodeAttestorServer interface {
	mustEmbedUnimplementedNodeAttestorServer()
}

func RegisterNodeAttestorServer(s grpc.ServiceRegistrar, srv NodeAttestorServer) {
	s.RegisterService(&NodeAttestor_ServiceDesc, srv)
}

func _NodeAttestor_Attest_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NodeAttestorServer).Attest(&nodeAttestorAttestServer{stream})
}

type NodeAttestor_AttestServer interface {
	Send(*AttestResponse) error
	Recv() (*AttestRequest, error)
	grpc.ServerStream
}

type nodeAttestorAttestServer struct {
	grpc.ServerStream
}

func (x *nodeAttestorAttestServer) Send(m *AttestResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nodeAttestorAttestServer) Recv() (*AttestRequest, error) {
	m := new(AttestRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NodeAttestor_ServiceDesc is the grpc.ServiceDesc for NodeAttestor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeAttestor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spire.plugin.server.nodeattestor.v1.NodeAttestor",
	HandlerType: (*NodeAttestorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Attest",
			Handler:       _NodeAttestor_Attest_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "spire/plugin/server/nodeattestor/v1/nodeattestor.proto",
}
