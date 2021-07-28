package noderesolver_test

import (
	"testing"

	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"github.com/spiffe/spire-plugin-sdk/plugintest"
	noderesolverv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/server/noderesolver/v1"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"github.com/spiffe/spire-plugin-sdk/templates/server/noderesolver"
)

func Test(t *testing.T) {
	plugin := new(noderesolver.Plugin)
	nrClient := new(noderesolverv1.NodeResolverPluginClient)
	configClient := new(configv1.ConfigServiceClient)

	// Serve the plugin in the background with the configured plugin and
	// service servers. The servers will be cleaned up when the test finishes.
	// TODO: Remove the config service server and client if no configuration
	// is required.
	// TODO: Provide host service server implementations if required by the
	// plugin.
	plugintest.ServeInBackground(t, plugintest.Config{
		PluginServer: noderesolverv1.NodeResolverPluginServer(plugin),
		PluginClient: nrClient,
		ServiceServers: []pluginsdk.ServiceServer{
			configv1.ConfigServiceServer(plugin),
		},
		ServiceClients: []pluginsdk.ServiceClient{
			configClient,
		},
	})

	// TODO: Invoke methods on the clients and assert the results
}
