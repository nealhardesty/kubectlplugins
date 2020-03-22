GOPATH?=$(HOME)/go
BUILD_OUT=build
.PHONY: build

KLOG_VERSION=v0.4.0
GNOSTIC_VERSION=v0.4.0

default: install


bootstrap:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

init: goget

goget: 
	go get ./cmd/...
	go mod tidy

build:
	mkdir -p $(BUILD_OUT)
	go build -o $(BUILD_OUT)/ ./cmd/...

vet:
	go vet ./cmd/...

install: 
	go install ./cmd/...

clean:
	go clean
	rm -rf build/
