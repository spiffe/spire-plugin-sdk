package keymanager

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl"
	"github.com/spiffe/spire-plugin-sdk/pluginmain"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
	keymanagerv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/server/keymanager/v1"
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

// Plugin implements the KeyManager plugin
type Plugin struct {
	// UnimplementedKeyManagerServer is embedded to satisfy gRPC
	keymanagerv1.UnimplementedKeyManagerServer

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

// GenerateKey implements the KeyManager GenerateKey RPC. Generates a new private key with the given ID.
// If a key already exists under that ID, it is overwritten and given a different fingerprint.
func (p *Plugin) GenerateKey(ctx context.Context, req *keymanagerv1.GenerateKeyRequest) (*keymanagerv1.GenerateKeyResponse, error) {
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

// GetPublicKey implements the KeyManager GetPublicKey RPC. Gets the public key information for the private key managed
// by the plugin with the given ID. If a key with the given ID does not exist, NOT_FOUND is returned.
func (p *Plugin) GetPublicKey(ctx context.Context, req *keymanagerv1.GetPublicKeyRequest) (*keymanagerv1.GetPublicKeyResponse, error) {
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

// GetPublicKeys implements the KeyManager GetPublicKeys RPC. Gets all public key information for the private keys
// managed by the plugin.
func (p *Plugin) GetPublicKeys(ctx context.Context, req *keymanagerv1.GetPublicKeysRequest) (*keymanagerv1.GetPublicKeysResponse, error) {
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

// SignData implements the KeyManager SignData RPC. Signs data with the private key identified by the given ID. If a key
// with the given ID does not exist, NOT_FOUND is returned. The response contains the signed data and the fingerprint of
// the key used to sign the data. See the PublicKey message for more details on the role of the fingerprint.
func (p *Plugin) SignData(ctx context.Context, req *keymanagerv1.SignDataRequest) (*keymanagerv1.SignDataResponse, error) {
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
		keymanagerv1.KeyManagerPluginServer(plugin),
		// TODO: Remove if no configuration is required
		configv1.ConfigServiceServer(plugin),
	)
}
