TAG = dev

PROJECT_ROOT := github.com/AmitKumarDas/storage-provisioner
PKG          := $(PROJECT_ROOT)/pkg
API_GROUPS   := ddp/v1alpha1

REGISTRY ?= quay.io/amitkumardas
IMG_NAME ?= storage-provisioner

BUILD_LDFLAGS = -X $(PROJECT_ROOT)/build.Hash=$(PACKAGE_VERSION)
GO_FLAGS = -gcflags '-N -l' -ldflags "$(BUILD_LDFLAGS)"

.PHONY: build
build: vendor generated_files $(IMG_NAME)

$(IMG_NAME):
	@echo "Making binary $@"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=off \
		go build -tags bins $(GO_FLAGS) -o $@ cmd/csi-attacher/main.go

.PHONY: unit-test
unit-test:
	pkgs="$$(go list ./... | grep -v '/test/integration/\|/examples/')" ; \
		go test -i $${pkgs} && \
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
	docker push $(REGISTRY)/$(IMG_NAME):$(TAG)

.PHONY: vendor
vendor:
	@echo "Vendor update ..."
	@dep ensure

.PHONY: generated_files
generated_files: deepcopy clientset lister informer

# deepcopy installs the deepcopy-gen at $GOPATH/bin
# Then make use of this installed binary to generate
# deepcopy
.PHONY: deepcopy
deepcopy:
	@go install ./vendor/k8s.io/code-generator/cmd/deepcopy-gen
	@echo "+ Generating deepcopy funcs for $(API_GROUPS)"
	@deepcopy-gen \
		--input-dirs $(PKG)/apis/$(API_GROUPS) \
		--output-file-base zz_generated.deepcopy \
		--go-header-file ./hack/custom-boilerplate.go.txt

# clienset installs the client-gen at $GOPATH/bin
# Then make use of this installed binary to generate
# clienset
.PHONY: clientset
clientset:
	@go install ./vendor/k8s.io/code-generator/cmd/client-gen
	@echo "+ Generating clientsets for $(API_GROUPS)"
	@client-gen \
		--fake-clientset=false \
		--input $(API_GROUPS) \
		--input-base $(PKG)/apis \
		--go-header-file ./hack/custom-boilerplate.go.txt \
		--clientset-path $(PKG)/client/generated/clientset

# lister installs the lister-gen at $GOPATH/bin
# Then make use of this installed binary to generate
# lister
.PHONY: lister
lister:
	@go install ./vendor/k8s.io/code-generator/cmd/lister-gen
	@echo "+ Generating lister for $(API_GROUPS)"
	@lister-gen \
		--input-dirs $(PKG)/apis/$(API_GROUPS) \
		--go-header-file ./hack/custom-boilerplate.go.txt \
		--output-package $(PKG)/client/generated/lister

# informer installs the informer-gen at $GOPATH/bin
# Then make use of this installed binary to generate
# informer
.PHONY: informer
informer:
	@go install ./vendor/k8s.io/code-generator/cmd/informer-gen
	@echo "+ Generating informer for $(API_GROUPS)"
	@informer-gen \
		--input-dirs $(PKG)/apis/$(API_GROUPS) \
		--output-package $(PKG)/client/generated/informer \
		--versioned-clientset-package $(PKG)/client/generated/clientset/internalclientset \
		--go-header-file ./hack/custom-boilerplate.go.txt \
		--listers-package $(PKG)/client/generated/lister