package plugintest_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	"github.com/spiffe/spire-plugin-sdk/private/proto/test"
)

func TestServe(t *testing.T) {
	log := new(bytes.Buffer)

	// Run in a subtest so that everything gets cleaned up before we
	// analyze the log.
	t.Run("", func(t *testing.T) {
		server := new(TestPlugin)
		pluginClient := new(test.SomePluginPluginClient)
		serviceClient := new(test.SomeServiceServiceClient)

		plugintest.ServeInBackground(t, plugintest.Config{
			PluginServer: test.SomePluginPluginServer(server),
			PluginClient: pluginClient,

			ServiceServers: []pluginsdk.ServiceServer{test.SomeServiceServiceServer(server)},
			ServiceClients: []pluginsdk.ServiceClient{serviceClient},

			HostServiceServers: []pluginsdk.ServiceServer{test.SomeHostServiceServiceServer(&someHostService{})},

			Logger: hclog.New(&hclog.LoggerOptions{
				DisableTime: true,
				Output:      log,
			}),
		})

		if server.log == nil {
			t.Fatal("logger should have been set")
		}
		if !server.hostServiceClient.IsInitialized() {
			t.Fatal("host service should have been initialized")
		}

		pluginResp, err := pluginClient.PluginEcho(context.Background(), &test.EchoRequest{In: "plugin-in"})
		if err != nil {
			t.Fatalf("PluginEcho failed unexpectedly with %v", err)
		}
		assertStringEqual(t, "plugin-in,plugin-out", pluginResp.Out)

		serviceResp, err := serviceClient.ServiceEcho(context.Background(), &test.EchoRequest{In: "service-in"})
		if err != nil {
			t.Fatalf("ServiceEcho failed unexpected with %v", err)
		}
		assertStringEqual(t, "service-in,service-out", serviceResp.Out)

		hostServiceResp, err := server.hostServiceClient.HostServiceEcho(context.Background(), &test.EchoRequest{In: "hostService-in"})
		if err != nil {
			t.Fatalf("HostServiceEcho failed unexpectedly with %v", err)
		}
		assertStringEqual(t, "hostService-in,hostService-out", hostServiceResp.Out)
	})

	assertStringEqual(t, "[INFO]  PLUGIN: in=plugin-in\n[INFO]  SERVICE: in=service-in\n", log.String())
}

func assertStringEqual(t *testing.T, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected string %q; got %q", expected, actual)
	}
}

type TestPlugin struct {
	test.UnimplementedSomePluginServer
	test.UnimplementedSomeServiceServer

	log               hclog.Logger
	hostServiceClient test.SomeHostServiceServiceClient
}

var _ pluginsdk.NeedsLogger = (*TestPlugin)(nil)
var _ pluginsdk.NeedsHostServices = (*TestPlugin)(nil)

func (p *TestPlugin) SetLogger(log hclog.Logger) {
	p.log = log
}

func (p *TestPlugin) BrokerHostServices(broker pluginsdk.ServiceBroker) error {
	if !broker.BrokerClient(&p.hostServiceClient) {
		return errors.New("host service was not available on broker")
	}
	return nil
}

func (p *TestPlugin) PluginEcho(_ context.Context, req *test.EchoRequest) (*test.EchoResponse, error) {
	p.log.Info("PLUGIN", "in", req.In)
	return &test.EchoResponse{Out: req.In + ",plugin-out"}, nil
}

func (p *TestPlugin) ServiceEcho(_ context.Context, req *test.EchoRequest) (*test.EchoResponse, error) {
	p.log.Info("SERVICE", "in", req.In)
	return &test.EchoResponse{Out: req.In + ",service-out"}, nil
}

type someHostService struct {
	test.UnimplementedSomeHostServiceServer
}

func (someHostService) HostServiceEcho(_ context.Context, req *test.EchoRequest) (*test.EchoResponse, error) {
	return &test.EchoResponse{Out: req.In + ",hostService-out"}, nil
}
