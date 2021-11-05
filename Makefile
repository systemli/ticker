GOCMD        = go
GOBUILD      = $(GOCMD) build
GOCLEAN      = $(GOCMD) clean
GOTEST       = $(GOCMD) test
GOGET        = $(GOCMD) get
GOX          = $(GOPATH)/bin/gox
GOGET        = $(GOCMD) get

GIT_VERSION  := $(shell git --no-pager describe --tags --always)
GIT_COMMIT   := $(shell git rev-parse --verify HEAD)

LD_FLAGS     = -X main.GitCommit=$(GIT_COMMIT) -X main.GitVersion=$(GIT_VERSION)

GOX_ARGS     = -output="$(BUILD_DIR)/{{.Dir}}-${GIT_VERSION}-{{.OS}}-{{.Arch}}" -osarch="linux/amd64 linux/arm linux/arm64 darwin/amd64 freebsd/amd64"

BUILD_DIR    = build
BINARY_NAME  = ticker

all: clean vet test build

build:
	$(GOBUILD) -ldflags "${LD_FLAGS}" -o $(BUILD_DIR)/$(BINARY_NAME) -v

test:
	${GOTEST} ./...

coverage:
	${GOTEST} -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/*

run: build
	cp config.yml.dist build/config.yml
	./$(BUILD_DIR)/$(BINARY_NAME) -config build/config.yml

release:
	${GOGET} -u github.com/mitchellh/gox
	${GOX} -ldflags "${LD_FLAGS}" ${GOX_ARGS}

docker:
	docker build --rm --force-rm --no-cache -t systemli/ticker .

.PHONY: all vet test coverage clean build run release docker
