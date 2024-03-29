// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package bundlepublisherv1

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

// BundlePublisherClient is the client API for BundlePublisher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BundlePublisherClient interface {
	// PublishBundle publishes the trust bundle that is in the request
	// to a store.
	PublishBundle(ctx context.Context, in *PublishBundleRequest, opts ...grpc.CallOption) (*PublishBundleResponse, error)
}

type bundlePublisherClient struct {
	cc grpc.ClientConnInterface
}

func NewBundlePublisherClient(cc grpc.ClientConnInterface) BundlePublisherClient {
	return &bundlePublisherClient{cc}
}

func (c *bundlePublisherClient) PublishBundle(ctx context.Context, in *PublishBundleRequest, opts ...grpc.CallOption) (*PublishBundleResponse, error) {
	out := new(PublishBundleResponse)
	err := c.cc.Invoke(ctx, "/spire.plugin.server.bundlepublisher.v1.BundlePublisher/PublishBundle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BundlePublisherServer is the server API for BundlePublisher service.
// All implementations must embed UnimplementedBundlePublisherServer
// for forward compatibility
type BundlePublisherServer interface {
	// PublishBundle publishes the trust bundle that is in the request
	// to a store.
	PublishBundle(context.Context, *PublishBundleRequest) (*PublishBundleResponse, error)
	mustEmbedUnimplementedBundlePublisherServer()
}

// UnimplementedBundlePublisherServer must be embedded to have forward compatible implementations.
type UnimplementedBundlePublisherServer struct {
}

func (UnimplementedBundlePublisherServer) PublishBundle(context.Context, *PublishBundleRequest) (*PublishBundleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishBundle not implemented")
}
func (UnimplementedBundlePublisherServer) mustEmbedUnimplementedBundlePublisherServer() {}

// UnsafeBundlePublisherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BundlePublisherServer will
// result in compilation errors.
type UnsafeBundlePublisherServer interface {
	mustEmbedUnimplementedBundlePublisherServer()
}

func RegisterBundlePublisherServer(s grpc.ServiceRegistrar, srv BundlePublisherServer) {
	s.RegisterService(&BundlePublisher_ServiceDesc, srv)
}

func _BundlePublisher_PublishBundle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishBundleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BundlePublisherServer).PublishBundle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.plugin.server.bundlepublisher.v1.BundlePublisher/PublishBundle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BundlePublisherServer).PublishBundle(ctx, req.(*PublishBundleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BundlePublisher_ServiceDesc is the grpc.ServiceDesc for BundlePublisher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BundlePublisher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spire.plugin.server.bundlepublisher.v1.BundlePublisher",
	HandlerType: (*BundlePublisherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishBundle",
			Handler:    _BundlePublisher_PublishBundle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spire/plugin/server/bundlepublisher/v1/bundlepublisher.proto",
}
