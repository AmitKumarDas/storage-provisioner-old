TAG = dev

PROJECT_ROOT := github.com/AmitKumarDas/storage-provisioner
PKG          := $(PROJECT_ROOT)/pkg
API_GROUPS   := ddp/v1alpha1

REGISTRY ?= quay.io/amitkumardas
IMG_NAME ?= storage-provisioner

BUILD_LDFLAGS = -X $(PROJECT_ROOT)/build.Hash=$(PACKAGE_VERSION)
GO_FLAGS = -gcflags '-N -l' -ldflags "$(BUILD_LDFLAGS)"

.PHONY: build
build: vendor generated_files unit-test $(IMG_NAME)

$(IMG_NAME):
	@echo "Making binary $@"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=off \
		go build -tags bins $(GO_FLAGS) -o $@ cmd/csi-attacher/main.go

.PHONY: unit-test
unit-test:
	@pkgs="$$(go list ./... | grep -v '/pkg/client/generated/')" ; \
		go test $${pkgs}

.PHONY: integration-test
integration-test:
	go test -i ./test/integration/...
	PATH="$(PWD)/hack/bin:$(PATH)" go test ./test/integration/... -v -timeout 5m -args -v=6

.PHONY: image
image: build
	@echo "Making image ..."
	docker build -t $(REGISTRY)/$(IMG_NAME):$(TAG) .

.PHONY: push
push: image
	@echo "Pushing image ..."
	@docker push $(REGISTRY)/$(IMG_NAME):$(TAG)

.PHONY: generated_files
generated_files: vendor
	@./hack/update-codegen.sh

.PHONY: vendor
vendor: go.mod go.sum
	@echo "Vendor update ..."
	@export GO111MODULE=on go mod download
	@export GO111MODULE=on go mod vendor