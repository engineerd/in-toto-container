PROJECT         := in-toto-container
ORG             := engineerd
BINDIR          := $(CURDIR)/bin
GOFLAGS         :=
LDFLAGS         := -w -s
TARGET_NAME     := itc

ifeq ($(OS),Windows_NT)
	TARGET = $(TARGET_NAME).exe
	SHELL  = cmd.exe
	CHECK  = where.exe
else
	TARGET = $(TARGET_NAME)
	SHELL  ?= bash
	CHECK  ?= which
endif

.PHONY: build
build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(TARGET) github.com/$(ORG)/$(PROJECT)/cmd/...

.PHONY: test
test:
	go test $(TESTFLAGS) ./... -v

.PHONY: e2e
e2e:
	$(BINDIR)/$(TARGET) verify --image radumatei/in-toto-container:v1 --layout testdata/intoto/demo.layout.template --layout-key testdata/intoto/alice.pub --links testdata/intoto/ --target testdata/intoto/foo.tar.gz

.PHONY: lint
lint:
	golangci-lint run --config ./golangci.yml

HAS_DEP          := $(shell $(CHECK) dep)
HAS_GOLANGCI     := $(shell $(CHECK) golangci-lint)
HAS_GOIMPORTS    := $(shell $(CHECK) goimports)
GOLANGCI_VERSION := v1.16.0


.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
ifndef HAS_GOLANGCI
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_VERSION)
endif
ifndef HAS_GOIMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
	dep check
