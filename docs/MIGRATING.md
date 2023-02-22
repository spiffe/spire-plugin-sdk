# Migrating Plugins

This document gives migration advice for porting existing plugins to the new
SDK interfaces.

## General Guidance

### Configuration

Pre-SDK plugin interfaces each contained their own `Configure` RPC. This RPC has
been replaced with a specific, versioned
[Config](/proto/spire/service/common/config) service.

See the relevant section in [AUTHORING](AUTHORING.md) for implementation guidance.

### GetPluginInfo

Pre-SDK plugin interfaces each included an unused `GetPluginInfo` RPC. This
RPC does not exist in the SDK plugin interfaces and can be removed.

### Errors

Pre-SDK plugins often added a prefix to all generated errors so that the
error messages logged by SPIRE would include some information about which
plugin originated the error. SPIRE now handles this prefixing automatically.
Accordingly, plugins should not prefix their own errors.

Migrated plugins _SHOULD_ return proper gRPC statuses when an error is
encountered, e.g.:

```
return status.Error(codes.InvalidArgument, "blah")
```

### Consuming Host Services

Plugins obtain Host Services clients by implementing the `BrokerHostServices`
function on the `pluginsdk.NeedsHostServices` interface. The function is passed
a broker that can be used to obtain the host service client. The broker
interface has changed slightly:

A pre-SDK implementation might look like this:

```
func (p *Plugin) BrokerHostServices(broker catalog.HostServiceBroker) error {
    has, err := broker.BrokerClient(foohostservicev1.FooHostServiceClient(&p.fooClient))
    if err != nil {
        return err
    }
    if !has {
        return errors.New("foo host service is required")
    }
    has, err = broker.BrokerClient(barhostservice.v1.BarHostServiceClient(&p.barClient))
    if err != nil {
        return err
    }
    if !has {
        p.log.Warn("Bar host service is not implemented")
    }
    return nil
}
```

And now looks like:

```
func (p *Plugin) BrokerHostServices(broker pluginsdk.ServiceBroker) error {
    if !broker.BrokerClient(&p.fooClient) {
        return errors.New("foo host service is required")
    }
    if !broker.BrokerClient(&p.barClient) {
        p.log.Warn("Bar host service is not implemented")
    }
    return nil
}
```

As before, plugin authors can decide if the lack of support for a specific host
service is an error or not. If the plugin returns an error from
`BrokerHostServices`, the plugin will fail to load.

### Common Types

Many of the pre-SDK plugin interfaces used types defined in `spire.common`.
Shared types are now defined in the SDK `spire.plugin.types` package.

## Plugin Specific Guidance

### Agent KeyManager

This interface has changed dramatically. The pre-SDK interface was a simple
interface for generation, storage, and retrieval of a single EC private key.

The SDK interface aligns with the same interface for the server, namely it
manages the creation and _use_ of multiple private key slots. Callers do not
use the private keys directly but rather do signing operations via a specific
RPC. As such, migrating to the new interface is _not_ trivial and will likely
require a rewrite of the plugin.

### Agent NodeAttestor

The `FetchAttestationData` RPC has been renamed to `AidAttestation`. The
request and response messages have been renamed to more closely align with the
flow. The response fields are now in a `oneof` to strongly convey the
difference in field requirements during the attestation flow.

### Agent WorkloadAttestor

The `Attest` RPC response returns selector values only. The selector type
is inferred by SPIRE from the name of the plugin.

### Server KeyManager

RSA 1024 bit key support has been dropped.

Keys managed by the KeyManager now have a required `fingerprint` component. The
fingerprint is returned when keys are created or retrieved by the `SignData`
RPC. It is used by SPIRE to detect if the key used in a signing operation has
changed from underneath SPIRE.

The `PSSOptions` message has been moved inside of the `SignDataRequest` message
to couple it to that operation.

### Server NodeAttestor

The `Attest` RPC request and response fields are now contained within `oneof`'s
to strongly convey the difference in field requirements in requests and
responses during the atestation flow. The attestation payload no longer needs
to include a type, since that is now inferred by SPIRE from the name of the
plugin. The selectors returned in the final response are selector values only.
The selector type is inferred by SPIRE from the name of the plugin.

### Server Notifier

No substantial changes outside of the migration to plugin SDK types.

### Server UpstreamAuthority

The `MintX509CA` RPC has been renamed to `MintX509CAAndSubscribe`. The response
uses the `X509Certificate` message from the plugin SDK types package instead of
raw ASN.1 bytes.

The `PublishJWTKey` RPC has been renamed to `PublishJWTKeyAndSubscribe`.
