package internal

import (
	"context"
	"errors"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/private"
	"google.golang.org/grpc"
)

// Serve serves the plugin with the given loggers and plugin/service servers and an optional test configuration.
func Serve(serverLogger, pluginLogger hclog.Logger, pluginServer pluginsdk.PluginServer, serviceServers []pluginsdk.ServiceServer, testConfig *goplugin.ServeTestConfig) {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: ServerHandshakeConfig(pluginServer),
		Plugins: map[string]goplugin.Plugin{
			pluginServer.Type(): newHCPlugin(serverLogger, pluginServer, serviceServers),
		},
		Logger:     pluginLogger,
		GRPCServer: goplugin.DefaultGRPCServer,
		Test:       testConfig,
	})
}

type hcServer struct {
	goplugin.NetRPCUnsupportedPlugin
	logger  hclog.Logger
	servers []pluginsdk.ServiceServer
}

func newHCPlugin(logger hclog.Logger, pluginServer pluginsdk.PluginServer, serviceServers []pluginsdk.ServiceServer) *hcServer {
	return &hcServer{
		logger:  logger,
		servers: append([]pluginsdk.ServiceServer{pluginServer}, serviceServers...),
	}
}

func (p *hcServer) GRPCServer(broker *goplugin.GRPCBroker, server *grpc.Server) (err error) {
	private.Register(server, p.servers, p.logger, &hcDialer{broker: broker})
	return nil
}

func (p *hcServer) GRPCClient(context.Context, *goplugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

type hcDialer struct {
	broker *goplugin.GRPCBroker
	conn   grpc.ClientConnInterface
}

func (d *hcDialer) DialHost(ctx context.Context) (grpc.ClientConnInterface, error) {
	if d.conn != nil {
		return d.conn, nil
	}

	conn, err := d.broker.Dial(private.HostServiceProviderID)
	if err != nil {
		return nil, err
	}
	d.conn = conn
	return conn, nil
}
