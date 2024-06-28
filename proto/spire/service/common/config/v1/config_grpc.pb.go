// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package configv1

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

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigClient interface {
	// Configure is called by SPIRE to configure the plugin with the plugin
	// specific configuration data and a set of SPIRE core configuration. It is
	// currently called when the plugin is first loaded after it has been
	// initialized. At a future point, it may be called to reconfigure the
	// plugin during runtime. Implementations should therefore expect that
	// calls to Configure can happen concurrently with other RPCs against the
	// plugin.
	Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*ConfigureResponse, error)
	// Validate is called by SPIRE to have the plugin determine if the specific
	// configuration data contains a valid configuration.
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*ConfigureResponse, error) {
	out := new(ConfigureResponse)
	err := c.cc.Invoke(ctx, "/spire.service.common.config.v1.Config/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/spire.service.common.config.v1.Config/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
// All implementations must embed UnimplementedConfigServer
// for forward compatibility
type ConfigServer interface {
	// Configure is called by SPIRE to configure the plugin with the plugin
	// specific configuration data and a set of SPIRE core configuration. It is
	// currently called when the plugin is first loaded after it has been
	// initialized. At a future point, it may be called to reconfigure the
	// plugin during runtime. Implementations should therefore expect that
	// calls to Configure can happen concurrently with other RPCs against the
	// plugin.
	Configure(context.Context, *ConfigureRequest) (*ConfigureResponse, error)
	// Validate is called by SPIRE to have the plugin determine if the specific
	// configuration data contains a valid configuration.
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
	mustEmbedUnimplementedConfigServer()
}

// UnimplementedConfigServer must be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (UnimplementedConfigServer) Configure(context.Context, *ConfigureRequest) (*ConfigureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}
func (UnimplementedConfigServer) Validate(context.Context, *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedConfigServer) mustEmbedUnimplementedConfigServer() {}

// UnsafeConfigServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServer will
// result in compilation errors.
type UnsafeConfigServer interface {
	mustEmbedUnimplementedConfigServer()
}

func RegisterConfigServer(s grpc.ServiceRegistrar, srv ConfigServer) {
	s.RegisterService(&Config_ServiceDesc, srv)
}

func _Config_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.service.common.config.v1.Config/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).Configure(ctx, req.(*ConfigureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spire.service.common.config.v1.Config/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Config_ServiceDesc is the grpc.ServiceDesc for Config service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Config_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spire.service.common.config.v1.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Configure",
			Handler:    _Config_Configure_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _Config_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spire/service/common/config/v1/config.proto",
}
