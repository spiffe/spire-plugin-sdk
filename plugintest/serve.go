package plugintest

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spiffe/spire-plugin-sdk/internal"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/private"
	"google.golang.org/grpc"
)

// Config is the test configuration for the plugin. It defines which plugin
// services, and host services are wired up.
type Config struct {
	// PluginServer is the plugin server implemented by the plugin.
	PluginServer pluginsdk.PluginServer

	// PluginClient is the plugin clients that the test needs have initialized
	// to facilitate testing the plugin implementations. This field is
	// optional.
	PluginClient pluginsdk.PluginClient

	// ServiceServers are the service servers implemented by the plugin. This
	// field is optional.
	ServiceServers []pluginsdk.ServiceServer

	// ServiceClients are the service clients that the test needs have
	// initialized to facilitate testing service implementations on the plugin.
	// This field is optional.
	ServiceClients []pluginsdk.ServiceClient

	// HostServiceServers are the host services that the test wants to offer to
	// the plugin. This field is optional.
	HostServiceServers []pluginsdk.ServiceServer

	// Logger is an optional logger for capturing log events from the plugin.
	// This field is optional.
	Logger hclog.Logger
}

// ServeInBackground serves the plugin in background with the provided
// configuration. The plugin will be unloaded when the test is over.
func ServeInBackground(t *testing.T, config Config) {
	wg := new(sync.WaitGroup)
	t.Cleanup(wg.Wait)

	switch {
	case config.PluginServer == nil:
		t.Fatal("PluginServer is required")
	case config.PluginClient == nil:
		t.Fatal("PluginClient is required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	closeCh := make(chan struct{})

	t.Cleanup(func() {
		cancel()
		<-closeCh
	})

	serverLogger := config.Logger
	if serverLogger == nil {
		serverLogger = hclog.NewNullLogger()
	}
	reattachConfigCh := make(chan *goplugin.ReattachConfig)

	if config.PluginServer == nil {
		t.Fatal("PluginServer must be provided")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		internal.Serve(serverLogger, hclog.NewNullLogger(), config.PluginServer, config.ServiceServers, &goplugin.ServeTestConfig{
			Context:          ctx,
			CloseCh:          closeCh,
			ReattachConfigCh: reattachConfigCh,
		})
	}()

	var reattachConfig *goplugin.ReattachConfig
	select {
	case reattachConfig = <-reattachConfigCh:
	case <-time.After(time.Second * 10):
		t.Fatal("timed out waiting for plugin to launch")
	}

	hcClient := &hcTestClient{wg: wg, hostServiceServers: config.HostServiceServers}

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: internal.ClientHandshakeConfig(config.PluginClient),
		Plugins: map[string]goplugin.Plugin{
			"TEST": hcClient,
		},
		Logger:           hclog.NewNullLogger(),
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		Reattach:         reattachConfig,
	})

	t.Cleanup(client.Kill)

	grpcClient, err := client.Client()
	if err != nil {
		t.Fatalf("failed to open gRPC client to plugin: %v", err)
	}
	t.Cleanup(func() { grpcClient.Close() })

	rawConn, err := grpcClient.Dispense("TEST")
	if err != nil {
		t.Fatalf("failed to dispense plugin client: %v", err)
	}

	conn, ok := rawConn.(*grpc.ClientConn)
	if !ok {
		// This is purely defensive; the hcTestClient GRPCClient method
		// returns a *grpc.ClientConn.
		t.Fatalf("expected %T; got %T", conn, rawConn)
	}
	t.Cleanup(func() { conn.Close() })

	grpcServiceNames, err := private.Init(ctx, conn, serverGRPCServiceNames(config.HostServiceServers))
	if err != nil {
		t.Fatalf("failed to initialize plugin: %v", err)
	}

	assertInitClient(t, conn, config.PluginClient, grpcServiceNames)
	for _, serviceClient := range config.ServiceClients {
		assertInitClient(t, conn, serviceClient, grpcServiceNames)
	}
}

func assertInitClient(t *testing.T, conn grpc.ClientConnInterface, client pluginsdk.ServiceClient, grpcServiceNames []string) {
	for _, grpcServiceName := range grpcServiceNames {
		if client.GRPCServiceName() == grpcServiceName {
			client.InitClient(conn)
			return
		}
	}
	t.Fatalf("%s was not provided by the plugin as expected", client.GRPCServiceName())
}

type hcTestClient struct {
	goplugin.NetRPCUnsupportedPlugin

	wg                 *sync.WaitGroup
	hostServiceServers []pluginsdk.ServiceServer
}

var _ goplugin.GRPCPlugin = (*hcTestClient)(nil)

func (p *hcTestClient) GRPCServer(b *goplugin.GRPCBroker, s *grpc.Server) error {
	return errors.New("not implemented host side")
}

func (p *hcTestClient) GRPCClient(ctx context.Context, b *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		b.AcceptAndServe(private.HostServiceProviderID, func(opts []grpc.ServerOption) *grpc.Server {
			s := grpc.NewServer(opts...)
			for _, hostServiceServer := range p.hostServiceServers {
				hostServiceServer.RegisterServer(s)
			}
			return s
		})
	}()
	return c, nil
}

func serverGRPCServiceNames(servers []pluginsdk.ServiceServer) []string {
	var grpcServiceNames []string
	for _, server := range servers {
		grpcServiceNames = append(grpcServiceNames, server.GRPCServiceName())
	}
	sort.Strings(grpcServiceNames)
	return grpcServiceNames
}
