package nodeattestor_test

import (
	"context"
	"testing"

	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	nodeattestorv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/server/nodeattestor/v1"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"github.com/spiffe/spire-plugin-sdk/templates/server/nodeattestor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	plugin := new(nodeattestor.Plugin)
	naClient := new(nodeattestorv1.NodeAttestorPluginClient)
	configClient := new(configv1.ConfigServiceClient)

	// Serve the plugin in the background with the configured plugin and
	// service servers. The servers will be cleaned up when the test finishes.
	// TODO: Remove the config service server and client if no configuration
	// is required.
	// TODO: Provide host service server implementations if required by the
	// plugin.
	plugintest.ServeInBackground(t, plugintest.Config{
		PluginServer: nodeattestorv1.NodeAttestorPluginServer(plugin),
		PluginClient: naClient,
		ServiceServers: []pluginsdk.ServiceServer{
			configv1.ConfigServiceServer(plugin),
		},
		ServiceClients: []pluginsdk.ServiceClient{
			configClient,
		},
	})

	ctx := context.Background()

	// TODO: Remove if no configuration is required.
	_, err := configClient.Configure(ctx, &configv1.ConfigureRequest{
		CoreConfiguration: &configv1.CoreConfiguration{TrustDomain: "example.org"},
		HclConfiguration:  `{}`,
	})
	assert.NoError(t, err)

	require.True(t, naClient.IsInitialized())

	// TODO: Make assertions using the desired plugin behavior.
	resp, err := naClient.Attest(context.Background())
	require.NoError(t, err)
	_, err = resp.Recv()
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
}
