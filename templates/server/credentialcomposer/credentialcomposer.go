package credentialcomposer

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl"
	"github.com/spiffe/spire-plugin-sdk/pluginmain"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	credentialcomposerv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/server/credentialcomposer/v1"
	configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// This compile time assertion ensures the plugin conforms properly to the
	// pluginsdk.NeedsLogger interface.
	// TODO: Remove if the plugin does not need the logger.
	_ pluginsdk.NeedsLogger = (*Plugin)(nil)

	// This compile time assertion ensures the plugin conforms properly to the
	// pluginsdk.NeedsHostServices interface.
	// TODO: Remove if the plugin does not need host services.
	_ pluginsdk.NeedsHostServices = (*Plugin)(nil)
)

// Config defines the configuration for the plugin.
// TODO: Add relevant configurables or remove if no configuration is required.
type Config struct {
}

// Plugin implements the CredentialComposer plugin
type Plugin struct {
	// UnimplementedCredentialComposerServer is embedded to satisfy gRPC
	credentialcomposerv1.UnimplementedCredentialComposerServer

	// UnimplementedConfigServer is embedded to satisfy gRPC
	// TODO: Remove if this plugin does not require configuration
	configv1.UnimplementedConfigServer

	// Configuration should be set atomically
	// TODO: Remove if this plugin does not require configuration
	configMtx sync.RWMutex
	config    *Config

	// The logger received from the framework via the SetLogger method
	// TODO: Remove if this plugin does not need the logger.
	logger hclog.Logger
}

// ComposeServerX509CA implements the CredentialComposer ComposeServerX509CA RPC. Composes the SPIRE Server X509 CA.
// The server will supply the default attributes it will apply to the CA. If the plugin returns an empty response or
// NOT_IMPLEMENTED, the server will apply the default attributes. Otherwise, the returned attributes are used.
// If a CA is produced that does not conform to the SPIFFE X509-SVID specification for signing certificates, it will be rejected.
func (p *Plugin) ComposeServerX509CA(ctx context.Context, req *credentialcomposerv1.ComposeServerX509CARequest) (*credentialcomposerv1.ComposeServerX509CAResponse, error) {
	config, err := p.getConfig()
	if err != nil {
		return nil, err
	}

	// TODO: Implement the RPC behavior. The following line silences compiler
	// warnings and can be removed once the configuration is referenced by the
	// implementation.
	config = config

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ComposeServerX509SVID implements the CredentialComposer ComposeServerX509SVID RPC. Composes the SPIRE Server X509-SVID.
// The server will supply the default attributes it will apply to the server X509-SVID. If the plugin returns an empty
// response or NOT_IMPLEMENTED, the server will apply the default attributes. Otherwise, the returned attributes are
// used. If an X509-SVID is produced that does not conform to the SPIFFE X509-SVID specification for leaf certificates,
// it will be rejected. This function cannot be used to modify the SPIFFE ID of the X509-SVID.
func (p *Plugin) ComposeServerX509SVID(ctx context.Context, req *credentialcomposerv1.ComposeServerX509SVIDRequest) (*credentialcomposerv1.ComposeServerX509SVIDResponse, error) {
	config, err := p.getConfig()
	if err != nil {
		return nil, err
	}

	// TODO: Implement the RPC behavior. The following line silences compiler
	// warnings and can be removed once the configuration is referenced by the
	// implementation.
	config = config

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ComposeAgentX509SVID implements the CredentialComposer ComposeAgentX509SVID RPC. Composes the SPIRE Agent X509-SVID.
// The server will supply the default attributes it will apply to the agent X509-SVID. If the plugin returns an empty
// response or NOT_IMPLEMENTED, the server will apply the default attributes. Otherwise, the returned attributes are used.
// If an X509-SVID is produced that does not conform to the SPIFFE X509-SVID specification for leaf certificates, it will
// be rejected. This function cannot be used to modify the SPIFFE ID of the X509-SVID.
func (p *Plugin) ComposeAgentX509SVID(ctx context.Context, req *credentialcomposerv1.ComposeAgentX509SVIDRequest) (*credentialcomposerv1.ComposeAgentX509SVIDResponse, error) {
	config, err := p.getConfig()
	if err != nil {
		return nil, err
	}

	// TODO: Implement the RPC behavior. The following line silences compiler
	// warnings and can be removed once the configuration is referenced by the
	// implementation.
	config = config

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ComposeWorkloadX509SVID implements the CredentialComposer ComposeWorkloadX509SVID RPC. Composes workload X509-SVIDs.
// The server will supply the default attributes it will apply to the workload X509-SVID. If the plugin returns an empty
// response or NOT_IMPLEMENTED, the server will apply the default attributes. Otherwise, the returned attributes are used.
// If an X509-SVID is produced that does not conform to the SPIFFE X509-SVID specification for leaf certificates, it will
// be rejected. This function cannot be used to modify the SPIFFE ID of the X509-SVID.
func (p *Plugin) ComposeWorkloadX509SVID(ctx context.Context, req *credentialcomposerv1.ComposeWorkloadX509SVIDRequest) (*credentialcomposerv1.ComposeWorkloadX509SVIDResponse, error) {
	config, err := p.getConfig()
	if err != nil {
		return nil, err
	}

	// TODO: Implement the RPC behavior. The following line silences compiler
	// warnings and can be removed once the configuration is referenced by the
	// implementation.
	config = config

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// ComposeWorkloadJWTSVID implements the CredentialComposer ComposeWorkloadJWTSVID RPC. Composes workload JWT-SVIDs.
// The server will supply the default attributes it will apply to the workload JWT-SVID. If the plugin returns an empty
// response or NOT_IMPLEMENTED, the server will apply the default attributes. Otherwise, the returned attributes are used.
// If a JWT-SVID is produced that does not conform to the SPIFFE JWT-SVID specification, it will be rejected.
// This function cannot be used to modify the SPIFFE ID of the JWT-SVID.
func (p *Plugin) ComposeWorkloadJWTSVID(ctx context.Context, req *credentialcomposerv1.ComposeWorkloadJWTSVIDRequest) (*credentialcomposerv1.ComposeWorkloadJWTSVIDResponse, error) {
	config, err := p.getConfig()
	if err != nil {
		return nil, err
	}

	// TODO: Implement the RPC behavior. The following line silences compiler
	// warnings and can be removed once the configuration is referenced by the
	// implementation.
	config = config

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// SetLogger is called by the framework when the plugin is loaded and provides
// the plugin with a logger wired up to SPIRE's logging facilities.
// TODO: Remove if the plugin does not need the logger.
func (p *Plugin) SetLogger(logger hclog.Logger) {
	p.logger = logger
}

// BrokerHostServices is called by the framework when the plugin is loaded to
// give the plugin a chance to obtain clients to SPIRE host services.
// TODO: Remove if the plugin does not need host services.
func (p *Plugin) BrokerHostServices(broker pluginsdk.ServiceBroker) error {
	// TODO: Use the broker to obtain host service clients
	return nil
}

// Configure configures the plugin. This is invoked by SPIRE when the plugin is
// first loaded. In the future, it may be invoked to reconfigure the plugin.
// As such, it should replace the previous configuration atomically.
// TODO: Remove if no configuration is required
func (p *Plugin) Configure(ctx context.Context, req *configv1.ConfigureRequest) (*configv1.ConfigureResponse, error) {
	config := new(Config)
	if err := hcl.Decode(config, req.HclConfiguration); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode configuration: %v", err)
	}

	// TODO: Validate configuration before setting/replacing existing
	// configuration

	p.setConfig(config)
	return &configv1.ConfigureResponse{}, nil
}

// setConfig replaces the configuration atomically under a write lock.
// TODO: Remove if no configuration is required
func (p *Plugin) setConfig(config *Config) {
	p.configMtx.Lock()
	p.config = config
	p.configMtx.Unlock()
}

// getConfig gets the configuration under a read lock.
// TODO: Remove if no configuration is required
func (p *Plugin) getConfig() (*Config, error) {
	p.configMtx.RLock()
	defer p.configMtx.RUnlock()
	if p.config == nil {
		return nil, status.Error(codes.FailedPrecondition, "not configured")
	}
	return p.config, nil
}

func main() {
	plugin := new(Plugin)
	// Serve the plugin. This function call will not return. If there is a
	// failure to serve, the process will exit with a non-zero exit code.
	pluginmain.Serve(
		credentialcomposerv1.CredentialComposerPluginServer(plugin),
		// TODO: Remove if no configuration is required
		configv1.ConfigServiceServer(plugin),
	)
}
