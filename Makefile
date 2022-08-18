DIR := ${CURDIR}

# There is no reason GOROOT should be set anymore. Unset it so it doesn't mess
# with our go toolchain detection/usage.
ifneq ($(GOROOT),)
    export GOROOT=
endif

E:=@
ifeq ($(V),1)
	E=
endif

cyan := $(shell which tput > /dev/null && tput setaf 6 2>/dev/null || echo "")
reset := $(shell which tput > /dev/null && tput sgr0 2>/dev/null || echo "")
bold  := $(shell which tput > /dev/null && tput bold 2>/dev/null || echo "")

.PHONY: default

default: generate

help:
	@echo "$(bold)Usage:$(reset) make $(cyan)<target>$(reset)"
	@echo "  $(cyan)generate$(reset)                              - generate gRPC and plugin interface code"
	@echo "  $(cyan)generate-check$(reset)                        - ensure generated code is up to date"
	@echo
	@echo "For verbose output set V=1"
	@echo "  for example: $(cyan)make V=1$(reset)"

# Used to force some rules to run every time
FORCE: ;

############################################################################
# Service definitions and other protos
############################################################################
plugin-protos := \
	private/proto/test/someplugin.proto \
	proto/spire/plugin/agent/keymanager/v1/keymanager.proto \
	proto/spire/plugin/agent/nodeattestor/v1/nodeattestor.proto \
	proto/spire/plugin/agent/svidstore/v1/svidstore.proto \
	proto/spire/plugin/agent/workloadattestor/v1/workloadattestor.proto \
	proto/spire/plugin/server/credentialcomposer/v1/credentialcomposer.proto \
	proto/spire/plugin/server/keymanager/v1/keymanager.proto \
	proto/spire/plugin/server/nodeattestor/v1/nodeattestor.proto \
	proto/spire/plugin/server/noderesolver/v1/noderesolver.proto \
	proto/spire/plugin/server/notifier/v1/notifier.proto \
	proto/spire/plugin/server/upstreamauthority/v1/upstreamauthority.proto \


service-protos := \
	private/proto/test/somehostservice.proto \
	private/proto/test/someservice.proto \
	proto/spire/hostservice/common/metrics/v1/metrics.proto \
	proto/spire/hostservice/server/agentstore/v1/agentstore.proto \
	proto/spire/hostservice/server/identityprovider/v1/identityprovider.proto \
	proto/spire/service/common/config/v1/config.proto \

grpc-protos := \
	internal/proto/spire/service/private/init/v1/init.proto \

protos := \
	private/proto/test/echo.proto \
	proto/spire/plugin/types/bundle.proto \
	proto/spire/plugin/types/jwtkey.proto \
	proto/spire/plugin/types/x509certificate.proto \

############################################################################
# OS/ARCH detection
############################################################################
os1=$(shell uname -s)
os2=
ifeq ($(os1),Darwin)
os1=darwin
os2=osx
else ifeq ($(os1),Linux)
os1=linux
os2=linux
else
$(error unsupported OS: $(os1))
endif

arch1=$(shell uname -m)
ifeq ($(arch1),x86_64)
arch2=amd64
else ifeq ($(arch1),aarch64)
arch2=arm64
else ifeq ($(arch1),arm64)
arch2=arm64
else
$(error unsupported ARCH: $(arch1))
endif

############################################################################
# Vars
############################################################################

build_dir := $(DIR)/.build/$(os1)-$(arch1)

go_version_full := $(shell cat .go-version)
go_version := $(go_version_full:.0=)
go_dir := $(build_dir)/go/$(go_version)
go_bin_dir := $(go_dir)/bin
go_url = https://storage.googleapis.com/golang/go$(go_version).$(os1)-$(arch2).tar.gz
go_path := PATH="$(go_bin_dir):$(PATH)"

golangci_lint_version = v1.27.0
golangci_lint_dir = $(build_dir)/golangci_lint/$(golangci_lint_version)
golangci_lint_bin = $(golangci_lint_dir)/golangci-lint

protoc_version = 3.20.1
ifeq ($(arch2),arm64)
protoc_url = https://github.com/protocolbuffers/protobuf/releases/download/v$(protoc_version)/protoc-$(protoc_version)-$(os2)-aarch_64.zip
else
protoc_url = https://github.com/protocolbuffers/protobuf/releases/download/v$(protoc_version)/protoc-$(protoc_version)-$(os2)-$(arch1).zip
endif
protoc_dir = $(build_dir)/protoc/$(protoc_version)
protoc_bin = $(protoc_dir)/bin/protoc

protoc_gen_go_version := $(shell grep google.golang.org/protobuf go.mod | awk '{print $$2}')
protoc_gen_go_base_dir := $(build_dir)/protoc-gen-go
protoc_gen_go_dir := $(protoc_gen_go_base_dir)/$(protoc_gen_go_version)-go$(go_version)
protoc_gen_go_bin := $(protoc_gen_go_dir)/protoc-gen-go

protoc_gen_go_grpc_version := v1.1.0
protoc_gen_go_grpc_base_dir := $(build_dir)/protoc-gen-go-grpc
protoc_gen_go_grpc_dir := $(protoc_gen_go_grpc_base_dir)/$(protoc_gen_go_grpc_version)-go$(go_version)
protoc_gen_go_grpc_bin := $(protoc_gen_go_grpc_dir)/protoc-gen-go-grpc

protoc_gen_go_spire_base_dir := $(build_dir)/protoc-gen-go-spire
protoc_gen_go_spire_dir := $(protoc_gen_go_spire_base_dir)/$(protoc_gen_go_spire_version)-go$(go_version)
protoc_gen_go_spire_bin := $(protoc_gen_go_spire_dir)/protoc-gen-go-spire

git_dirty := $(shell git --no-pager status -s)

#############################################################################
# Utility functions and targets
#############################################################################

.PHONY: git-clean-check

git-clean-check:
ifneq ($(git_dirty),)
	git --no-pager diff
	@echo "Git repository is dirty!"
	@false
else
	@echo "Git repository is clean."
endif

#############################################################################
# Test Targets
#############################################################################

.PHONY: test

test: | go-check
ifneq ($(COVERPROFILE),)
	$(E)$(go_path) go test $(go_flags) -covermode=atomic -coverprofile="$(COVERPROFILE)" ./...
else
	$(E)$(go_path) go test $(go_flags) ./...
endif

#############################################################################
# Code Generation
#############################################################################

.PHONY: generate generate-plugin-protos generate-service-protos generate-grpc-protos generate-protos generate-check

generate: \
	generate-plugin-protos \
	generate-service-protos \
	generate-grpc-protos \
	generate-protos

generate-plugin-protos: \
	$(plugin-protos:.proto=.pb.go) \
	$(plugin-protos:.proto=_grpc.pb.go) \
	$(plugin-protos:.proto=_spire_plugin.pb.go)

generate-service-protos: \
	$(service-protos:.proto=.pb.go) \
	$(service-protos:.proto=_grpc.pb.go) \
	$(service-protos:.proto=_spire_service.pb.go)

generate-grpc-protos: \
	$(grpc-protos:.proto=.pb.go) \
	$(grpc-protos:.proto=_grpc.pb.go)

generate-protos: \
	$(protos:.proto=.pb.go)

parentdir = $(patsubst %/,%,$(dir $(1)))
protobase-rec = $(if $(patsubst .,,$(1)), \
		$(if $(filter-out proto,$(notdir $(1))),$(call protobase-rec,$(call parentdir,$(1)),$(2)),$(1)), \
		$(error could not find proto base of $(2)) \
	)
protobase = $(call protobase-rec,$(1),$(1))

%_spire_plugin.pb.go: %.proto $(protoc_bin) $(protoc_gen_go_spire_bin) FORCE
	@echo "generating $@..."
	$(E) PATH="$(protoc_gen_go_spire_dir):$(PATH)" $(protoc_bin) \
		-I $(call protobase,$@) \
		--go-spire_out=. \
		--go-spire_opt=module=github.com/spiffe/spire-plugin-sdk \
		--go-spire_opt=mode=plugin \
		$<

%_spire_service.pb.go: %.proto $(protoc_bin) $(protoc_gen_go_spire_bin) FORCE
	@echo "generating $@..."
	$(E) PATH="$(protoc_gen_go_spire_dir):$(PATH)" $(protoc_bin) \
		-I $(call protobase,$@) \
		--go-spire_out=. \
		--go-spire_opt=module=github.com/spiffe/spire-plugin-sdk \
		--go-spire_opt=mode=service \
		$<

%_grpc.pb.go: %.proto $(protoc_bin) $(protoc_gen_go_grpc_bin) FORCE
	@echo "generating $@..."
	$(E) PATH="$(protoc_gen_go_grpc_dir):$(PATH)" $(protoc_bin) \
		-I $(call protobase,$@) \
		--go-grpc_out=. --go-grpc_opt=module=github.com/spiffe/spire-plugin-sdk \
		$<

%.pb.go: %.proto $(protoc_bin) $(protoc_gen_go_bin) FORCE
	@echo "generating $@..."
	$(E) PATH="$(protoc_gen_go_dir):$(PATH)" $(protoc_bin) \
		-I $(call protobase,$@) \
		--go_out=. --go_opt=module=github.com/spiffe/spire-plugin-sdk \
		$<

generate-check:
ifneq ($(git_dirty),)
	$(error generate-check must be invoked on a clean repository)
endif
	$(E)find . -type f -name "*.proto" -exec touch {} \;
	@echo "Compiling protocol buffers..."
	$(E)$(MAKE) generate
	@echo "Ensuring git repository is clean..."
	$(E)$(MAKE) git-clean-check

#############################################################################
# Toolchain
#############################################################################

# go-check checks to see if there is a version of Go available matching the
# required version. The build cache is preferred. If not available, it is
# downloaded into the build cache. Any rule needing to invoke tools in the go
# toolchain should depend on this rule and then prepend $(go_bin_dir) to their
# path before invoking go or use $(go_path) go which already has the path prepended.
# Note that some tools (e.g. anything that uses golang.org/x/tools/go/packages)
# execute on the go binary and also need the right path in order to locate the
# correct go binary.
go-check:
ifneq (go$(go_version), $(shell $(go_path) go version 2>/dev/null | cut -f3 -d' '))
	@echo "Installing go$(go_version)..."
	$(E)rm -rf $(dir $(go_dir))
	$(E)mkdir -p $(go_dir)
	$(E)curl -sSfL $(go_url) | tar xz -C $(go_dir) --strip-components=1
endif

go-bin-path: go-check
	@echo "$(go_bin_dir):${PATH}"

$(protoc_bin):
	@echo "Installing protoc $(protoc_version)..."
	$(E)rm -rf $(dir $(protoc_dir))
	$(E)mkdir -p $(protoc_dir)
	$(E)curl -sSfL $(protoc_url) -o $(build_dir)/tmp.zip; unzip -q -d $(protoc_dir) $(build_dir)/tmp.zip; rm $(build_dir)/tmp.zip

$(protoc_gen_go_bin): | go-check
	@echo "Installing protoc-gen-go $(protoc_gen_go_version)..."
	$(E)rm -rf $(protoc_gen_go_base_dir)
	$(E)mkdir -p $(protoc_gen_go_dir)
	$(E)$(go_path) go build -o $(protoc_gen_go_bin) google.golang.org/protobuf/cmd/protoc-gen-go

$(protoc_gen_go_grpc_bin): | go-check
	@echo "Installing protoc-gen-go-grpc $(protoc_gen_go_grpc_version)..."
	$(E)rm -rf $(protoc_gen_go_grpc_base_dir)
	$(E)mkdir -p $(protoc_gen_go_grpc_dir)
	$(E)echo "module tools" > $(protoc_gen_go_grpc_dir)/go.mod
	$(E)cd $(protoc_gen_go_grpc_dir) && GOBIN=$(protoc_gen_go_grpc_dir) $(go_path) go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(protoc_gen_go_grpc_version)

$(protoc_gen_go_spire_bin): $(wildcard ./cmd/protoc-gen-go-spire/*) | go-check
	@echo "Installing protoc-gen-go-spire..."
	$(E)rm -rf $(protoc_gen_go_spire_base_dir)
	$(E)mkdir -p $(protoc_gen_go_spire_dir)
	$(E)go build -o $(protoc_gen_go_spire_bin) ./cmd/protoc-gen-go-spire
