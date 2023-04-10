# Authoring Plugins

This document gives guidance for authoring plugins.

SPIRE plugins implement one and only one plugin _type_ (e.g. KeyManager). They
also implement zero or more services. Below is a list of plugin types, alongside templates that can be used as a base
for authoring plugins.


## Templates

### Agent

| Plugin           | Description                                           | Template                                    |
|------------------|-------------------------------------------------------|---------------------------------------------|
| KeyManager       | Manages private keys and performs signing operations. | [link](../templates/agent/keymanager)       |
| NodeAttestor     | Performs the agent side of the node attestation flow. | [link](../templates/agent/nodeattestor)     |
| SVIDStore        | Stores workload X509-SVIDs to arbitrary destinations. | [link](../templates/agent/svidstore)        |
| WorkloadAttestor | Attests workloads and provides selectors.             | [link](../templates/agent/workloadattestor) |

### Server

| Plugin            | Description                                            | Template                                      |
|-------------------|--------------------------------------------------------|-----------------------------------------------|
| KeyManager        | Manages private keys and performs signing operations.  | [link](../templates/server/keymanager)        |
| NodeAttestor      | Performs the server side of the node attestation flow. | [link](../templates/server/nodeattestor)      |
| NodeResolver      | Provides additional selectors for attested nodes.      | [link](../templates/server/noderesolver)      |
| Notifier          | Notifies external systems of certain SPIRE events.     | [link](../templates/server/notifier)          |
| UpstreamAuthority | Plugs SPIRE into an upstream PKI.                      | [link](../templates/server/upstreamauthority) |


## Configuration

Most plugins require some form of configuration. SPIRE supports passing plugin
configuration data to plugins when they are loaded. This configuration data is
provided to SPIRE in the `plugin_data` section of the plugin declaration in the
server or agent configuration, e.g.:

```
plugins {
    UpstreamAuthority "disk" {
        plugin_data {
            key_file_path = "some.key"
            cert_file_path = "some.crt"
        }
    }
```

In order to receive this configuration data, as well as other core
configurables, a plugin implements the [Config](/proto/spire/service/common/config) service:

Implementing this service is **optional**; plugins which do not require
additional configuration are free to not implement it. However, if SPIRE
receives configuration data for a plugin and that plugin does NOT implement
a configuration service, SPIRE will fail to load the plugin.

To implement this service

- Embed the `UnimplementedConfigServer` struct:

```
import configv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/service/common/config/v1"

...

type Plugin struct {
    configv1.UnimplementedConfigServer
}
```

- Implement the `Configure` RPC:

```
type Config struct {
    // ... HCL fields ...
}

func (*Plugin) Configure(ctx context.Context, req *configv1.ConfigureRequest) (*configv1.ConfigureResponse, error) {
    config := new(Config)
    if err := hcl.Decode(config, req.HclConfiguration); err != nil {
        return status.Errorf(codes.InvalidArgument, "failed to decode configuration: %v", err)
    }

    ...
    return &configv1.ConfigureResponse{}, nil
}
```

- Advertise that the implementation implements the configuration service:

```
func main() {
    plugin := new(Plugin)
    pluginmain.Serve(
        keymanagerv1.KeyManagerPluginServer(plugin),
        configv1.ConfigServiceServer(plugin), // <-- add the Config service server implementation
    )
}
```

## Errors

Plugins _SHOULD_ return proper gRPC statuses when an error is encountered.
SPIRE will automatically prefix all errors that originate from a plugin with
the name and type of the plugin. For example, if a `keymanager` plugin with the
name `foo` returns an error like so:

```
return status.Error(codes.InvalidArgument, "blah")
```

Then SPIRE will prefix the error as it is received, producing the equivalent
of:

```
return status.Error(codes.InvalidArgument, "keymanager(foo): blah")
```

This helps identify the source of errors.

## Logging

SPIRE provides plugins a logger that is wired up to the SPIRE logging
facilities. Plugin authors should use this logger instead of their own in order
to maintain a consistent logging experience.

To access the logger, plugin implementations implement the `pluginsdk.NeedsLogger`
interface, like so:

```
func (p *Plugin) SetLogger(logger hclog.Logger) {
}
```

The passed in logger can be stored by the plugin for later use.

## Consuming Host Services

Plugins obtain Host Services clients by implementing the `BrokerHostServices`
function on the `pluginsdk.NeedsHostServices` interface. The function is passed
a broker that can be used to obtain the host service client.

For example:

```
type Plugin struct {
    fooClient foohostservicev1.FooServiceClient
    barClient barhostservicev1.BarServiceClient
    // ... other fields...
}


func (p *Plugin) BrokerHostServices(broker pluginsdk.ServiceBroker) error {
    if !broker.BrokerClient(&p.fooClient) {
        return errors.New("foo host service is required")
    }
    if !broker.BrokerClient(&p.barClient) {
        p.log.Warn("Bar host service is not implemented")
    }
    return nil
}

func (p *Plugin) SomeMethod() {
    // Since the bar client is optional, it should only be used if it was initialized
    if p.barClient.IsInitialized() {
        ...
    }
}
```

Plugin authors can decide if the lack of support for a specific host service is
an error or not. If the plugin returns an error from BrokerHostServices, the
plugin will fail to load.

## Cleanup

Plugins are separate processes and are terminated when the plugin is unloaded.
However, it may be desirable to perform some graceful cleanup operations.

To facilitate this, if plugin/service implementations implement the io.Closer
interface, then the `Close` method will be invoked before the plugin is
unloaded. No other RPCs will be invoked at any time during or after the `Close`
method is called. Errors returned from `Close` are simply logged and will not
impact any runtime behavior of SPIRE Server.

Implementations of `Close` should avoid long running or blocking behavior.
SPIRE may employ deadlines on the operation and could terminate the plugin
before the cleanup is fully completed if plugin implementations ignore this
advice.

## Unit Testing

The [plugintest](https://pkg.go.dev/github.com/spiffe/spire-plugin-sdk/plugintest) 
package can be used to conveniently test plugin implementations. The test framework
loads the plugin in the background and hosts the specified plugin/service/hostservice
servers. It initializes clients that can be used to invoke RPCs and test functionality.

See the package docs for more information.

## Running

The [pluginmain](https://pkg.go.dev/github.com/spiffe/spire-plugin-sdk/pluginmain) package
is used to run the plugin. It takes care of setting up all the plugin facilities and
wiring up the logger and hostservices.

See the package docs for more information.
