SHELL = bash
.DEFAULT_GOAL = build

$(shell mkdir -p bin)
export GOBIN = $(realpath bin)
export PATH := $(PATH):$(GOBIN)
export OS   := $(shell if [ "$(shell uname)" = "Darwin" ]; then echo "darwin"; else echo "linux"; fi)
export ARCH := $(shell if [ "$(shell uname -m)" = "x86_64" ]; then echo "amd64"; else echo "arm64"; fi)

# Extracts the version number for a given dependency found in go.mod.
# Makes the test setup be in sync with what the operator itself uses.
extract-version = $(shell cat go.mod | grep $(1) | awk '{$$1=$$1};1' | cut -d' ' -f2 | sed 's/^v//')

#### TOOLS ####
TOOLS_DIR                          := $(PWD)/.tools
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.20.0
CONTROLLER_GEN_VERSION             := $(call extract-version,sigs.k8s.io/controller-tools)
CHAINSAW_VERSION                   := $(call extract-version,github.com/kyverno/chainsaw)

#### VARS ####
KUBE_CONTEXT               ?= kind-$(KIND_CLUSTER_NAME)
KUBERNETES_VERSION          = 1.29.0
KIND_IMAGE                 ?= kindest/node:v$(KUBERNETES_VERSION)
KIND_CLUSTER_NAME          ?= bekk

.PHONY: generate
generate:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v${CONTROLLER_GEN_VERSION}
	./bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd
	go generate ./...

.PHONY: build
build: generate
	go build \
	-tags osusergo,netgo \
	-trimpath \
	-ldflags="-s -w" \
	-o ./bin/faggrupperator \
	./cmd

.PHONY: run-local
run-local: build install-operator
	kubectl --context ${KUBE_CONTEXT} apply -f config/ --recursive
	./bin/faggrupperator

.PHONY: setup-local
setup-local: kind-cluster install-operator
	@echo "Cluster $(KUBE_CONTEXT) is setup"


#### KIND ####

.PHONY: kind-cluster check-kind
check-kind:
	@which kind >/dev/null || (echo "kind not installed, please install it to proceed"; exit 1)

.PHONY: kind-cluster
kind-cluster: check-kind
	@echo Create kind cluster... >&2
	@kind create cluster --image $(KIND_IMAGE) --name ${KIND_CLUSTER_NAME}


.PHONY: install-operator
install-operator: generate
	@kubectl create namespace bekk-system --context $(KUBE_CONTEXT) || true
	@kubectl create namespace demo --context $(KUBE_CONTEXT) || true
	@kubectl apply -f config/ --recursive --context $(KUBE_CONTEXT)

.PHONY: install-test-tools
install-test-tools:
	go install github.com/kyverno/chainsaw@v${CHAINSAW_VERSION}

#### TESTS ####
.PHONY: test-single
test-single: install-test-tools install-operator
	@./bin/chainsaw test --kube-context $(KUBE_CONTEXT) --config tests/config.yaml --test-dir $(dir) && \
    echo "Test succeeded" || (echo "Test failed" && exit 1)

.PHONY: test
test: install-test-tools install-operator
	@./bin/chainsaw test --kube-context $(KUBE_CONTEXT) --config tests/config.yaml --test-dir tests/ && \
    echo "Test succeeded" || (echo "Test failed" && exit 1)

.PHONY: run-test
run-test: build
	@echo "Starting faggrupperator in background..."
	@LOG_FILE=$$(mktemp -t faggrupperator-test.XXXXXXX); \
	./bin/faggrupperator -e error > "$$LOG_FILE" 2>&1 & \
	PID=$$!; \
	echo "faggrupperator PID: $$PID"; \
	echo "Log redirected to file: $$LOG_FILE"; \
	( \
		if [ -z "$(TEST_DIR)" ]; then \
			$(MAKE) test; \
		else \
			$(MAKE) test-single dir=$(TEST_DIR); \
		fi; \
	) && \
	(echo "Stopping faggrupperator (PID $$PID)..." && kill $$PID)  || (echo "Test or faggrupperator failed. Stopping faggrupperator (PID $$PID)" && kill $$PID && exit 1)
