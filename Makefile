GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOX=$(GOPATH)/bin/gox

GOX_ARGS = -output="$(BUILD_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64 linux/386 linux/arm linux/arm64 darwin/amd64 freebsd/amd64 freebsd/386 windows/386 windows/amd64"

BUILD_DIR=build
BINARY_NAME=ticker

all: clean vet test build

build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

vet:
	$(GOVET) ./...

test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/*

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)
	cp config.yml.dist build/config.yml
	./$(BUILD_DIR)/$(BINARY_NAME)

release:
	$(GOGET) -u github.com/mitchellh/gox
	$(GOX) $(GOX_ARGS)
