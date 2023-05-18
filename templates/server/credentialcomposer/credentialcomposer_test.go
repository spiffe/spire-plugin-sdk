package credentialcomposer_test

import (
	"context"
	"testing"

	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	credentialcomposerv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/server/credentialcomposer/v1"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"github.com/spiffe/spire-plugin-sdk/templates/server/credentialcomposer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	plugin := new(credentialcomposer.Plugin)
	pluginClient := new(credentialcomposerv1.CredentialComposerPluginClient)
	configClient := new(configv1.ConfigServiceClient)

	// Serve the plugin in the background with the configured plugin and
	// service servers. The servers will be cleaned up when the test finishes.
	// TODO: Remove the config service server and client if no configuration
	// is required.
	// TODO: Provide host service server implementations if required by the
	// plugin.
	plugintest.ServeInBackground(t, plugintest.Config{
		PluginServer: credentialcomposerv1.CredentialComposerPluginServer(plugin),
		PluginClient: pluginClient,
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

	require.True(t, pluginClient.IsInitialized())
	// TODO: Make assertions using the desired plugin behavior.
	_, err = pluginClient.ComposeServerX509CA(ctx, &credentialcomposerv1.ComposeServerX509CARequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
	_, err = pluginClient.ComposeServerX509SVID(ctx, &credentialcomposerv1.ComposeServerX509SVIDRequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
	_, err = pluginClient.ComposeAgentX509SVID(ctx, &credentialcomposerv1.ComposeAgentX509SVIDRequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
	_, err = pluginClient.ComposeWorkloadX509SVID(ctx, &credentialcomposerv1.ComposeWorkloadX509SVIDRequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
	_, err = pluginClient.ComposeWorkloadJWTSVID(ctx, &credentialcomposerv1.ComposeWorkloadJWTSVIDRequest{})
	assert.EqualError(t, err, "rpc error: code = Unimplemented desc = not implemented")
}
