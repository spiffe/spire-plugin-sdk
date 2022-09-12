package private

import (
	"context"

	initv1 "github.com/spiffe/spire-plugin-sdk/internal/proto/spire/service/private/init/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Deinit deinitializes the plugin. It should only be called right before the
// host unloads the plugin and will not be invoking any other plugin or service
// RPCs.
func Deinit(ctx context.Context, conn grpc.ClientConnInterface) error {
	client := initv1.NewInitClient(conn)
	_, err := client.Deinit(ctx, &initv1.DeinitRequest{})
	switch status.Code(err) {
	case codes.OK, codes.Unimplemented:
		return nil
	default:
		return err
	}
}
