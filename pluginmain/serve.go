package pluginmain

import (
	"github.com/spiffe/spire-plugin-sdk/internal"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
)

// Serve serves the plugin using the given plugin/service servers. It does
// not return. It is intended to be called from main(). For example:
//
//	func main() {
//	    plugin := new(Plugin)
//	    pluginmain.Serve(
//	         keymanagerv1.KeyManagerPluginServer(plugin),
//	         configv1.ConfigServiceServer(plugin),
//	    )
//	}
func Serve(pluginServer pluginsdk.PluginServer, serviceServers ...pluginsdk.ServiceServer) {
	logger := internal.NewLogger()
	internal.Serve(logger, logger, pluginServer, serviceServers, nil)
}
