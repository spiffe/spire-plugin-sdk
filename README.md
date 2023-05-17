# SPIRE Plugin SDK

This repository contains the service definitions, code generated stubs, and
infrastructure for running and testing [SPIRE](https://github.com/spiffe/spire) plugins.

## Overview

SPIRE supports a rich plugin system. Plugins can either be built in, or
external, to SPIRE. External plugins are separate processes and use
[go-plugin](https://github.com/hashicorp/go-plugin) under the covers.

SPIRE communicates with plugins over gRPC. As such, the various interfaces are defined via gRPC service definitions.

There are three types of interfaces:

| Type         | Description
| ------------ | --------------------------------------------------------------|
| Plugin       | The primary plugin interface. A plugin only implements only one plugin interface. |
| Service      | An auxiliary service interface. These are generic facilities consumed by SPIRE. An example is the common [Config](proto/spire/service/common/config) service. A plugin implements zero or more service interfaces. |
| Host Service | A service provided by SPIRE and optionally consumed by plugins. |

## Plugins

### Agent

| Plugin | Versions | Description | Template    |
| ------ | -------- | ----------- | ----------- |
| KeyManager       | [v1](proto/spire/plugin/agent/keymanager/v1/keymanager.proto)                | Manages private keys and performs signing operations.  | [link](templates/agent/keymanager)         |
| NodeAttestor     | [v1](proto/spire/plugin/agent/nodeattestor/v1/nodeattestor.proto)            | Performs the agent side of the node attestation flow.  | [link](templates/agent/nodeattestor)       |
| SVIDStore        | [v1](proto/spire/plugin/agent/svidstore/v1/svidstore.proto)                  | Stores workload X509-SVIDs to arbitrary destinations.  | [link](templates/agent/svidstore)          |
| WorkloadAttestor | [v1](proto/spire/plugin/agent/workloadattestor/v1/workloadattestor.proto)    | Attests workloads and provides selectors.              | [link](templates/agent/workloadattestor)   |

### Server

| Plugin | Versions  | Description | Template    |
| ------ | --------  | ----------- | ----------- |
| BundlePublisher | [v1](proto/spire/plugin/server/bundlepublisher/v1/bundlepublisher.proto) | Publishes a trust bundle to a store.               | [link](templates/server/bundlepublisher) |
| CredentialComposer | [v1](proto/spire/plugin/server/credentialcomposer/v1/credentialcomposer.proto) | Allows customization of SVID and CA attributes.        | [link](templates/server/credentialcomposer) |
| KeyManager         | [v1](proto/spire/plugin/server/keymanager/v1/keymanager.proto)                 | Manages private keys and performs signing operations.  | [link](templates/server/keymanager)         |
| NodeAttestor       | [v1](proto/spire/plugin/server/nodeattestor/v1/nodeattestor.proto)             | Performs the server side of the node attestation flow. | [link](templates/server/nodeattestor)       |
| Notifier           | [v1](proto/spire/plugin/server/notifier/v1/notifier.proto)                     | Notifies external systems of certain SPIRE events.     | [link](templates/server/notifier)           |
| UpstreamAuthority  | [v1](proto/spire/plugin/server/upstreamauthority/v1/upstreamauthority.proto)   | Plugs SPIRE into an upstream PKI.                      | [link](templates/server/upstreamauthority)  |


## Services

### Common

| Service | Versions | Description |
| ------- | -------- | ----------- |
| Config | [v1](proto/spire/service/common/config/v1/config.proto) | Used by SPIRE to configure the plugin. |


## Host Services

### Common

| Host Service | Versions | Description |
| ------------ | -------- | ----------- |
| Metrics | [v1](proto/spire/hostservice/common/metrics/v1/metrics.proto) | Provides metrics facilities. |


### Server

| Host Service | Versions | Description |
| ------------ | -------- | ----------- |
| IdentityProvider | [v1](proto/spire/hostservice/server/identityprovider/v1/identityprovider.proto) | Provides an identity and bundle information. |
| AgentStore       | [v1](proto/spire/hostservice/server/agentstore/v1/agentstore.proto)             | Provides information about attested agents.  |


## Authoring Plugins

For guidance in authoring a plugin, see [AUTHORING](/docs/AUTHORING.md).

## Migrating Pre-SDK Plugins

To migrate existing pre-SDK plugins, see [MIGRATING](/docs/MIGRATING.md).

## Versioning

This repository is tagged along with SPIRE releases with the same name, even if
there are no changes to the APIs between SPIRE versions. This allows consumers
to always pick a tag that matches up with their deployment. Even so, SPIRE
maintains API compatibility between SPIRE versions. SPIRE will clearly indicate
in the [CHANGELOG](https://github.com/spiffe/spire/blob/main/CHANGELOG) when
APIs are deprecated and issue warnings at runtime when they are used well in
advance of any removal.

## Contributing

This repository follows the same governance and contribution guidelines as the
[SPIRE](https://github.com/spiffe/spire) project.

For specifics on getting started, see [CONTRIBUTING](/docs/CONTRIBUTING.md).

Please open [Issues](https://github.com/spiffe/spire/issues) to request features or file bugs.
