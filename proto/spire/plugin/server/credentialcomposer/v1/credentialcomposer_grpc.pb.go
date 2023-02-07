// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package credentialcomposerv1

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

// CredentialComposerClient is the client API for CredentialComposer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CredentialComposerClient interface {
	// Composes the SPIRE Server X509 CA. The server will supply the default
	// attributes it will apply to the CA. If the plugin returns an empty
	// response or NOT_IMPLEMENTED, the server will apply the default
	// attributes. Otherwise the returned attributes are used. If a CA is
	// produced that does not conform to the SPIFFE X509-SVID specification for
	// signing certificates, it will be rejected.
	ComposeServerX509CA(ctx context.Context, in *ComposeServerX509CARequest, opts ...grpc.CallOption) (*ComposeServerX509CAResponse, error)
	// Composes the SPIRE Server X509-SVID. The server will supply the default
	// attributes it will apply to the server X509-SVID. If the plugin returns
	// an empty response or NOT_IMPLEMENTED, the server will apply the default
	// attributes. Otherwise the returned attributes are used. If an X509-SVID
	// is produced that does not conform to the SPIFFE X509-SVID specification
	// for leaf certificates, it will be rejected. This function cannot be used
	// to modify the SPIFFE ID of the X509-SVID.
	ComposeServerX509SVID(ctx context.Context, in *ComposeServerX509SVIDRequest, opts ...grpc.CallOption) (*ComposeServerX509SVIDResponse, error)
	// Composes the SPIRE Agent X509-SVID. The server will supply the default
	// attributes it will apply to the agent X509-SVID. If the plugin returns
	// an empty response or NOT_IMPLEMENTED, the server will apply the default
	// attributes. Otherwise the returned attributes are used. If an X509-SVID
	// is produced that does not conform to the SPIFFE X509-SVID specification
	// for leaf certificates, it will be rejected. This function cannot be used
	// to modify the SPIFFE ID of the X509-SVID.
	ComposeAgentX509SVID(ctx context.Context, in *ComposeAgentX509SVIDRequest, opts ...grpc.CallOption) (*ComposeAgentX509SVIDResponse, error)
	// Composes workload X509-SVIDs. The server will supply the default
	// attributes it will apply to the workload X509-SVID. If the plugin
	// returns an empty response or NOT_IMPLEMENTED, the server will apply the
	// default attributes. Otherwise the returned attributes are used. If an
	// X509-SVID is produced that does not conform to the SPIFFE X509-SVID
	// specification for leaf certificates, it will be rejected. This function
	// cannot be used to modify the SPIFFE ID of the X509-SVID.
	ComposeWorkloadX509SVID(ctx context.Context, in *ComposeWorkloadX509SVIDRequest, opts ...grpc.CallOption) (*ComposeWorkloadX509SVIDResponse, error)
	// Composes workload JWT-SVIDs. The server will supply the default
	// attributes it will apply to the workload JWT-SVID. If the plugin
	// returns an empty response or NOT_IMPLEMENTED, the server will apply the
	// default attributes. Otherwise the returned attributes are used. If a
	// JWT-SVID is produced that does not conform to the SPIFFE JWT-SVID
	// specification, it will be rejected. This function cannot be used to
	// modify the SPIFFE ID of the JWT-SVID.
	ComposeWorkloadJWTSVID(ctx context.Context, in *ComposeWorkloadJWTSVIDRequest, opts ...grpc.CallOption) (*ComposeWorkloadJWTSVIDResponse, error)
}

type credentialComposerClient struct {
	cc grpc.ClientConnInterface
}

func NewCredentialComposerClient(cc grpc.ClientConnInterface) CredentialComposerClient {
	return &credentialComposerClient{cc}
}

func (c *credentialComposerClient) ComposeServerX509CA(ctx context.Context, in *ComposeServerX509CARequest, opts ...grpc.CallOption) (*ComposeServerX509CAResponse, error) {
	out := new(ComposeServerX509CAResponse)
	err := c.cc.Invoke(ctx, "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeServerX509CA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credentialComposerClient) ComposeServerX509SVID(ctx context.Context, in *ComposeServerX509SVIDRequest, opts ...grpc.CallOption) (*ComposeServerX509SVIDResponse, error) {
	out := new(ComposeServerX509SVIDResponse)
	err := c.cc.Invoke(ctx, "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeServerX509SVID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credentialComposerClient) ComposeAgentX509SVID(ctx context.Context, in *ComposeAgentX509SVIDRequest, opts ...grpc.CallOption) (*ComposeAgentX509SVIDResponse, error) {
	out := new(ComposeAgentX509SVIDResponse)
	err := c.cc.Invoke(ctx, "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeAgentX509SVID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credentialComposerClient) ComposeWorkloadX509SVID(ctx context.Context, in *ComposeWorkloadX509SVIDRequest, opts ...grpc.CallOption) (*ComposeWorkloadX509SVIDResponse, error) {
	out := new(ComposeWorkloadX509SVIDResponse)
	err := c.cc.Invoke(ctx, "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeWorkloadX509SVID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *credentialComposerClient) ComposeWorkloadJWTSVID(ctx context.Context, in *ComposeWorkloadJWTSVIDRequest, opts ...grpc.CallOption) (*ComposeWorkloadJWTSVIDResponse, error) {
	out := new(ComposeWorkloadJWTSVIDResponse)
	err := c.cc.Invoke(ctx, "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeWorkloadJWTSVID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CredentialComposerServer is the server API for CredentialComposer service.
// All implementations must embed UnimplementedCredentialComposerServer
// for forward compatibility
type CredentialComposerServer interface {
	// Composes the SPIRE Server X509 CA. The server will supply the default
	// attributes it will apply to the CA. If the plugin returns an empty
	// response or NOT_IMPLEMENTED, the server will apply the default
	// attributes. Otherwise the returned attributes are used. If a CA is
	// produced that does not conform to the SPIFFE X509-SVID specification for
	// signing certificates, it will be rejected.
	ComposeServerX509CA(context.Context, *ComposeServerX509CARequest) (*ComposeServerX509CAResponse, error)
	// Composes the SPIRE Server X509-SVID. The server will supply the default
	// attributes it will apply to the server X509-SVID. If the plugin returns
	// an empty response or NOT_IMPLEMENTED, the server will apply the default
	// attributes. Otherwise the returned attributes are used. If an X509-SVID
	// is produced that does not conform to the SPIFFE X509-SVID specification
	// for leaf certificates, it will be rejected. This function cannot be used
	// to modify the SPIFFE ID of the X509-SVID.
	ComposeServerX509SVID(context.Context, *ComposeServerX509SVIDRequest) (*ComposeServerX509SVIDResponse, error)
	// Composes the SPIRE Agent X509-SVID. The server will supply the default
	// attributes it will apply to the agent X509-SVID. If the plugin returns
	// an empty response or NOT_IMPLEMENTED, the server will apply the default
	// attributes. Otherwise the returned attributes are used. If an X509-SVID
	// is produced that does not conform to the SPIFFE X509-SVID specification
	// for leaf certificates, it will be rejected. This function cannot be used
	// to modify the SPIFFE ID of the X509-SVID.
	ComposeAgentX509SVID(context.Context, *ComposeAgentX509SVIDRequest) (*ComposeAgentX509SVIDResponse, error)
	// Composes workload X509-SVIDs. The server will supply the default
	// attributes it will apply to the workload X509-SVID. If the plugin
	// returns an empty response or NOT_IMPLEMENTED, the server will apply the
	// default attributes. Otherwise the returned attributes are used. If an
	// X509-SVID is produced that does not conform to the SPIFFE X509-SVID
	// specification for leaf certificates, it will be rejected. This function
	// cannot be used to modify the SPIFFE ID of the X509-SVID.
	ComposeWorkloadX509SVID(context.Context, *ComposeWorkloadX509SVIDRequest) (*ComposeWorkloadX509SVIDResponse, error)
	// Composes workload JWT-SVIDs. The server will supply the default
	// attributes it will apply to the workload JWT-SVID. If the plugin
	// returns an empty response or NOT_IMPLEMENTED, the server will apply the
	// default attributes. Otherwise the returned attributes are used. If a
	// JWT-SVID is produced that does not conform to the SPIFFE JWT-SVID
	// specification, it will be rejected. This function cannot be used to
	// modify the SPIFFE ID of the JWT-SVID.
	ComposeWorkloadJWTSVID(context.Context, *ComposeWorkloadJWTSVIDRequest) (*ComposeWorkloadJWTSVIDResponse, error)
	mustEmbedUnimplementedCredentialComposerServer()
}

// UnimplementedCredentialComposerServer must be embedded to have forward compatible implementations.
type UnimplementedCredentialComposerServer struct {
}

func (UnimplementedCredentialComposerServer) ComposeServerX509CA(context.Context, *ComposeServerX509CARequest) (*ComposeServerX509CAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComposeServerX509CA not implemented")
}
func (UnimplementedCredentialComposerServer) ComposeServerX509SVID(context.Context, *ComposeServerX509SVIDRequest) (*ComposeServerX509SVIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComposeServerX509SVID not implemented")
}
func (UnimplementedCredentialComposerServer) ComposeAgentX509SVID(context.Context, *ComposeAgentX509SVIDRequest) (*ComposeAgentX509SVIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComposeAgentX509SVID not implemented")
}
func (UnimplementedCredentialComposerServer) ComposeWorkloadX509SVID(context.Context, *ComposeWorkloadX509SVIDRequest) (*ComposeWorkloadX509SVIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComposeWorkloadX509SVID not implemented")
}
func (UnimplementedCredentialComposerServer) ComposeWorkloadJWTSVID(context.Context, *ComposeWorkloadJWTSVIDRequest) (*ComposeWorkloadJWTSVIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ComposeWorkloadJWTSVID not implemented")
}
func (UnimplementedCredentialComposerServer) mustEmbedUnimplementedCredentialComposerServer() {}

// UnsafeCredentialComposerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CredentialComposerServer will
// result in compilation errors.
type UnsafeCredentialComposerServer interface {
	mustEmbedUnimplementedCredentialComposerServer()
}

func RegisterCredentialComposerServer(s grpc.ServiceRegistrar, srv CredentialComposerServer) {
	s.RegisterService(&CredentialComposer_ServiceDesc, srv)
}

func _CredentialComposer_ComposeServerX509CA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComposeServerX509CARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialComposerServer).ComposeServerX509CA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeServerX509CA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialComposerServer).ComposeServerX509CA(ctx, req.(*ComposeServerX509CARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredentialComposer_ComposeServerX509SVID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComposeServerX509SVIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialComposerServer).ComposeServerX509SVID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeServerX509SVID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialComposerServer).ComposeServerX509SVID(ctx, req.(*ComposeServerX509SVIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredentialComposer_ComposeAgentX509SVID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComposeAgentX509SVIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialComposerServer).ComposeAgentX509SVID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeAgentX509SVID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialComposerServer).ComposeAgentX509SVID(ctx, req.(*ComposeAgentX509SVIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredentialComposer_ComposeWorkloadX509SVID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComposeWorkloadX509SVIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialComposerServer).ComposeWorkloadX509SVID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeWorkloadX509SVID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialComposerServer).ComposeWorkloadX509SVID(ctx, req.(*ComposeWorkloadX509SVIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CredentialComposer_ComposeWorkloadJWTSVID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComposeWorkloadJWTSVIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialComposerServer).ComposeWorkloadJWTSVID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.plugin.server.credentialcomposer.v1.CredentialComposer/ComposeWorkloadJWTSVID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialComposerServer).ComposeWorkloadJWTSVID(ctx, req.(*ComposeWorkloadJWTSVIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CredentialComposer_ServiceDesc is the grpc.ServiceDesc for CredentialComposer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CredentialComposer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spire.plugin.server.credentialcomposer.v1.CredentialComposer",
	HandlerType: (*CredentialComposerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ComposeServerX509CA",
			Handler:    _CredentialComposer_ComposeServerX509CA_Handler,
		},
		{
			MethodName: "ComposeServerX509SVID",
			Handler:    _CredentialComposer_ComposeServerX509SVID_Handler,
		},
		{
			MethodName: "ComposeAgentX509SVID",
			Handler:    _CredentialComposer_ComposeAgentX509SVID_Handler,
		},
		{
			MethodName: "ComposeWorkloadX509SVID",
			Handler:    _CredentialComposer_ComposeWorkloadX509SVID_Handler,
		},
		{
			MethodName: "ComposeWorkloadJWTSVID",
			Handler:    _CredentialComposer_ComposeWorkloadJWTSVID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spire/plugin/server/credentialcomposer/v1/credentialcomposer.proto",
}