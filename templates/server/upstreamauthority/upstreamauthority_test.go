package upstreamauthority_test

import (
	"context"
	"testing"

	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	upstreamauthorityv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/server/upstreamauthority/v1"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"github.com/spiffe/spire-plugin-sdk/templates/server/upstreamauthority"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	plugin := new(upstreamauthority.Plugin)
	uaClient := new(upstreamauthorityv1.UpstreamAuthorityPluginClient)
	configClient := new(configv1.ConfigServiceClient)

	// Serve the plugin in the background with the configured plugin and
	// service servers. The servers will be cleaned up when the test finishes.
	// TODO: Remove the config service server and client if no configuration
	// is required.
	// TODO: Provide host service server implementations if required by the
	// plugin.
	plugintest.ServeInBackground(t, plugintest.Config{
		PluginServer: upstreamauthorityv1.UpstreamAuthorityPluginServer(plugin),
		PluginClient: uaClient,
		ServiceServers: []pluginsdk.ServiceServer{
			configv1.ConfigServiceServer(plugin),
		},
		ServiceClients: []pluginsdk.ServiceClient{
			configClient,
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: Remove if no configuration is required.
	_, err := configClient.Configure(ctx, &configv1.ConfigureRequest{
		CoreConfiguration: &configv1.CoreConfiguration{TrustDomain: "example.org"},
		HclConfiguration:  `{}`,
	})
	assert.NoError(t, err)

	require.True(t, uaClient.IsInitialized())
	// TODO: Make assertions using the desired plugin behavior.
	mintStream, err := uaClient.MintX509CAAndSubscribe(ctx, &upstreamauthorityv1.MintX509CARequest{})
	require.NoError(t, err)
	_, err = mintStream.Recv()
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
	publishStream, err := uaClient.PublishJWTKeyAndSubscribe(ctx, &upstreamauthorityv1.PublishJWTKeyRequest{})
	require.NoError(t, err)
	_, err = publishStream.Recv()
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
}
