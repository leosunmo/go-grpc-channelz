BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT_VERSION := v1.17.1
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
PROTOC_VERSION := 3.9.1
RELEASE_OS :=
PROTOC_DIR := .tmp/protoc-$(PROTOC_VERSION)
PROTOC_BIN := $(PROTOC_DIR)/bin/protoc

ifeq ($(OS),Windows_NT)
	echo "Windows not supported yet, sorry...."
	exit 1
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		RELEASE_OS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		RELEASE_OS = osx
	endif
endif


all: test lint

tidy:
	go mod tidy -v

build: protoc
	go build ./...

test: build
	go test -cover -race ./...

test-coverage:
	go test ./... -race -coverprofile=coverage.txt && go tool cover -html=coverage.txt

ci-test: build
	go test -race $$(go list ./...) -v -coverprofile coverage.txt -covermode=atomic

setup: setup-git-hooks

setup-git-hooks:
	git config core.hooksPath .githooks

lint: $(GOLANGCI_LINT)
	$(GOPATH)/bin/golangci-lint run --fast --enable-all -D gochecknoglobals

$(GOLANGCI_LINT):
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi


$(PROTOC_BIN):
	@echo "Installing unzip (if required)"
	@which unzip || apt-get update || sudo apt-get update
	@which unzip || apt-get install unzip || sudo apt-get install unzip
	@echo Installing protoc
	rm -rf $(PROTOC_DIR)
	mkdir -p $(PROTOC_DIR)
	cd $(PROTOC_DIR) &&\
		curl -OL https://github.com/google/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(RELEASE_OS)-x86_64.zip &&\
		unzip protoc-$(PROTOC_VERSION)-$(RELEASE_OS)-x86_64.zip
	chmod +x $(PROTOC_BIN)
	@echo "Installing protoc-gen-go (if required)"
	@which protoc-gen-go > /dev/null || GO111MODULE=on go get -u github.com/golang/protobuf/protoc-gen-go

run-demo-server:
	go run internal/demo/server/main/main.go

run-demo-client:
	go run internal/demo/client/main/main.go

protoc: $(PROTOC_BIN)
	mkdir -p internal/generated/service
	$(PROTOC_BIN) --proto_path=internal/proto/ --go_out=plugins=grpc:internal/generated/service internal/proto/*.proto

