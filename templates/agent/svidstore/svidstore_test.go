package svidstore_test

import (
	"context"
	"testing"

	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	svidstorev1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/agent/svidstore/v1"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"github.com/spiffe/spire-plugin-sdk/templates/agent/svidstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	plugin := new(svidstore.Plugin)
	ssClient := new(svidstorev1.SVIDStorePluginClient)
	configClient := new(configv1.ConfigServiceClient)

	// Serve the plugin in the background with the configured plugin and
	// service servers. The servers will be cleaned up when the test finishes.
	// TODO: Remove the config service server and client if no configuration
	// is required.
	// TODO: Provide host service server implementations if required by the
	// plugin.
	plugintest.ServeInBackground(t, plugintest.Config{
		PluginServer: svidstorev1.SVIDStorePluginServer(plugin),
		PluginClient: ssClient,
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

	require.True(t, ssClient.IsInitialized())
	// TODO: Make assertions using the desired plugin behavior.
	_, err = ssClient.PutX509SVID(ctx, &svidstorev1.PutX509SVIDRequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
	_, err = ssClient.DeleteX509SVID(ctx, &svidstorev1.DeleteX509SVIDRequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
}
