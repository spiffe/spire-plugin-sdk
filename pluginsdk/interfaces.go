package pluginsdk

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

// PluginServer is implemented by plugin server implementations.
type PluginServer interface {
	ServiceServer

	// Type returns the type of plugin (e.g. "KeyManager")
	Type() string
}

// PluginClient is implemented by plugin client implementations.
type PluginClient interface {
	ServiceClient

	// Type returns the type of plugin (e.g. "KeyManager")
	Type() string
}

// ServiceServer is implemented by service/hostservice server implementations.
type ServiceServer interface {
	// GRPCServiceName returns the full gRPC service name (e.g.
	// "spire.plugin.server.keymanager.v1.KeyManager")
	GRPCServiceName() string

	// RegisterServer registers the server implementation with the given gRPC server.
	// It returns the implementation that was registered.
	RegisterServer(server *grpc.Server) interface{}
}

// ServiceClient is implemented by service/hostservice client implementations
type ServiceClient interface {
	// GRPCServiceName returns the full gRPC service name (e.g.
	// "spire.plugin.server.keymanager.v1.KeyManager")
	GRPCServiceName() string

	// InitClient initializes the client using the given gRPC client
	// connection. It returns the client implementation that was initialized.
	InitClient(conn grpc.ClientConnInterface) interface{}
}

// NeedsLogger is an interface implemented by server implementations that
// a logger. The provided logger is wired up to SPIRE and logs emitted with
// the logger will show up in SPIRE logs.
type NeedsLogger interface {
	SetLogger(logger hclog.Logger)
}

// NeedsHostServices is an interface implemented by plugin/service
// server implementations that need SPIRE host services.
type NeedsHostServices interface {
	// BrokerHostServices is invoked by the plugin loader and provides a broker
	// that can be used by plugins/services to initialize clients to host
	// services provided by SPIRE. If an error is returned, plugin loading will
	// fail. This gives server implementations control over whether or not the
	// absence of a particular host service is a catastrophic failure.
	BrokerHostServices(ServiceBroker) error
}

// ServiceBroker is used to obtain clients to SPIRE host services.
type ServiceBroker interface {
	// BrokerClient initializes the passed in host service client. If the
	// host service is not available in SPIRE, the host service client will
	// remain uninitialized and the function will return false.
	BrokerClient(ServiceClient) bool
}
