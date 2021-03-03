
default: build

build: manifests format

ifeq ($(VERSION),)
VERSION=$(shell git describe --tags  --long)-$(shell date +"%Y%m%d%H%M%S")
endif

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

CRD_OPTIONS ?= "crd:trivialVersions=true"

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	mkdir -p config/crd/db
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/db/v1/..."
	$(CONTROLLER_GEN) $(CRD_OPTIONS) paths="./api/db/..." output:crd:artifacts:config=config/crd/db output:none

format: .bin/yq
	 for file in $$(find config -iname *.yaml); do $(YQ) r -P $$file > $$file.tmp; mv $$file.tmp $$file; done

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

OS   = $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH = $(shell uname -m | sed 's/x86_64/amd64/')

.bin/yq:
	mkdir -p .bin
	curl -sSLo .bin/yq https://github.com/mikefarah/yq/releases/download/3.4.1/yq_$(OS)_$(ARCH) && chmod +x .bin/yq

YQ = $(realpath ./.bin/yq)