package private

import (
	"context"

	initv1 "github.com/spiffe/spire-plugin-sdk/internal/proto/spire/service/private/init/v1"
	"google.golang.org/grpc"
)

// Init initializes the plugin and advertises the given host service names to
// the plugin for brokering. The list of service names implemented by the
// plugin are returned. This function is only intended to be used internally
// and by SPIRE. See /private/README.md.
func Init(ctx context.Context, conn grpc.ClientConnInterface, hostServiceNames []string) (pluginServiceNames []string, err error) {
	client := initv1.NewInitClient(conn)
	resp, err := client.Init(ctx, &initv1.InitRequest{
		HostServiceNames: hostServiceNames,
	})
	if err != nil {
		return nil, err
	}
	return resp.PluginServiceNames, nil
}
