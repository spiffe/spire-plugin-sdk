package private

import (
	"context"

	"github.com/hashicorp/go-hclog"
	initv1 "github.com/spiffe/spire-plugin-sdk/internal/proto/spire/service/private/init/v1"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	"google.golang.org/grpc"
)

// HostDialer is a generic interface for dialing the host (i.e. SPIRE). This
// interface is only intended to be used internally and by SPIRE. See
// /private/README.md.
type HostDialer interface {
	DialHost(ctx context.Context) (grpc.ClientConnInterface, error)
}

// Register registers the given servers with the gRPC server. The given and
// logger will be used when the plugins are initialized. This function is only
// intended to be used internally and by SPIRE. See /private/README.md.
func Register(s *grpc.Server, servers []pluginsdk.ServiceServer, logger hclog.Logger, dialer HostDialer) {
	var names []string
	var impls []interface{}
	for _, server := range servers {
		names = append(names, server.GRPCServiceName())
		impls = append(impls, server.RegisterServer(s))
	}

	initv1.RegisterInitServer(s, &initService{
		logger: logger,
		names:  names,
		impls:  impls,
		dialer: dialer,
	})
}

type initService struct {
	initv1.UnimplementedInitServer

	logger hclog.Logger
	names  []string
	impls  []interface{}
	dialer HostDialer
}

func (s *initService) Init(ctx context.Context, req *initv1.InitRequest) (*initv1.InitResponse, error) {
	initted := map[interface{}]struct{}{}
	for _, impl := range s.impls {
		// Wire up the logger and host service broker. Since the same
		// implementation might back more than one server, only initialize
		// once.
		if _, ok := initted[impl]; ok {
			continue
		}
		initted[impl] = struct{}{}

		if impl, ok := impl.(pluginsdk.NeedsLogger); ok {
			impl.SetLogger(s.logger)
		}

		if impl, ok := impl.(pluginsdk.NeedsHostServices); ok {
			conn, err := s.dialer.DialHost(ctx)
			if err != nil {
				return nil, err
			}
			broker := serviceBroker{conn: conn, hostServiceNames: req.HostServiceNames}
			if err := impl.BrokerHostServices(broker); err != nil {
				s.logger.Error("Plugin failed brokering host services", "error", err)
				return nil, err
			}
		}
	}

	return &initv1.InitResponse{
		PluginServiceNames: s.names,
	}, nil
}

type serviceBroker struct {
	conn             grpc.ClientConnInterface
	hostServiceNames []string
}

func (b serviceBroker) BrokerClient(client pluginsdk.ServiceClient) bool {
	wants := client.GRPCServiceName()
	for _, has := range b.hostServiceNames {
		if wants == has {
			client.InitClient(b.conn)
			return true
		}
	}
	return false
}
