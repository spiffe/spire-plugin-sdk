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

## Consuming Changes in SPIRE

SPIRE's main branch depends on a pseudo-version of this repository (see
https://golang.org/ref/mod#pseudo-versions).

While a new change in this repository is under development, you can add a
temporary `replace` directive to the SPIRE `go.mod` to allow you to consume the
changes.  Care must be taken to not push the `replace` directive change up to
SPIRE.

Once those changes have been merged and you are ready to consume them from
SPIRE, run `go get github.com/spiffe/spire-plugin-sdk@<commit hash>` in the SPIRE
repository. This will update `go.mod` in SPIRE to use the latest pseudo version
with that commit.

When cutting a SPIRE release, this repository is tagged with the SPIRE
release version. The release branch in SPIRE is updated to depend explicitly
on that version (i.e. `go get github.com/spiffe/spire-plugin-sdk@<version>`).

Relying on a pseudo versions means that this repository only needs tags
for the offically released versions, while still allowing SPIRE to work with
unreleased changes during development.
