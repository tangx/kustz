
BINARY:=kustz

GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)
VERSION=v$(shell cat .version)

WORKSPACE ?= ./cmd/kustz

tidy:
	go mod tidy

build:
	go build -o out/$(BINARY)-$(VERSION)-$(GOOS)-$(GOARCH) $(WORKSPACE)

install:
	go install $(WORKSPACE)
	
	
build.x:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(MAKE) build
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(MAKE) build
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 $(MAKE) build
	CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 $(MAKE) build

clean:
	rm -rf out


test.deployment:
	make test TARGET=Test_KustzDeployment
test.kustomize:
	make test TARGET=Test_KustzKustomize


test:
	go test -timeout 30s -run ^$(TARGET) github.com/tangx/kustz/pkg/kustz -v -count=1
	
