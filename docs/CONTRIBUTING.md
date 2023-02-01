# Contribution Guidelines

This document outlines the contribution guidelines for the SPIRE Plugin SDK.

This project follows the contribution and governance guidelines from the SPIFFE
project (see
[CONTRIBUTING](https://github.com/spiffe/spiffe/blob/main/CONTRIBUTING.md)
and [GOVERNANCE](https://github.com/spiffe/spiffe/blob/main/GOVERNANCE.md)).

## Prerequisites

The [Makefile](/Makefile) in the project is set up to download required
dependencies for code generation.

### Updating Dependency Versions

The [Makefile](/Makefile) uses internal variables or inspects [go.mod](/go.mod)
to determine the versions of various tools in the toolchain. See the
[Makefile](/Makefile) for specifics.

## Generating Service Definitions

To (re)generate service definitions do the following:

```sh
$ make
```

If you are adding a new .proto file, you first need to update the `Makefile`
and add the .proto file to the relevant variables.

## Opening PRs

All PRs should target the `next` branch. The `next` branch is a staging area
for all features under development but not ready for release in an official
version of SPIRE.

Changes are cherry-picked into `main` from the `next` branch ahead of an
official SPIRE release. The commits in `main` are tagged with the supporting
SPIRE version.

## Consuming Changes in SPIRE

While a new change in this repository is under development, you can use [Go
Workspaces](https://go.dev/ref/mod#workspaces) to allow SPIRE to consume the
changes before they are merged into this repository.

SPIRE's main branch depends on a pseudo-version of this repository based on the
`next` branch (see https://golang.org/ref/mod#pseudo-versions). Once changes
have been merged into the `next` branch, the pseudo-version dependency in the
SPIRE repository can be updated by running `go get
github.com/spiffe/spire-plugin-sdk@next` from the SPIRE repository.

Relying on a pseudo versions means that this repository only needs tags
for the offically released versions, while still allowing SPIRE to work with
unreleased changes during development.
