package internal

import (
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spiffe/spire-plugin-sdk/pluginsdk"
)

// ServerHandshakeConfig returns the handshake configuration for the given
// server implementation.
func ServerHandshakeConfig(pluginServer pluginsdk.PluginServer) goplugin.HandshakeConfig {
	return goplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   pluginServer.Type(),
		MagicCookieValue: pluginServer.Type(),
	}
}

// ClientHandshakeConfig returns the handshake configuration for the given
// client implementation.
func ClientHandshakeConfig(pluginClient pluginsdk.PluginClient) goplugin.HandshakeConfig {
	return goplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   pluginClient.Type(),
		MagicCookieValue: pluginClient.Type(),
	}
}
