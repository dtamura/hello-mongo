NAME			:= hello-mongo
VERSION			:= v0.0.1
REVISION		:= $(shell git rev-parse --short HEAD)
PACKAGE_PATH	:= github.com/dtamura/${NAME}
SRCS    		:= $(shell find . -type f -name '*.go')
LDFLAGS			:= "-s -w -X \"$(PACKAGE_PATH)/server.Name=$(NAME)\" -X \"$(PACKAGE_PATH)/server.Version=$(VERSION)\" -X \"$(PACKAGE_PATH)/server.Revision=$(REVISION)\" -extldflags \"-static\""

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

LINT := $(GOBIN)/golint

$(LINT): ; @go get github.com/golang/lint/golint
$(GODOC): ; @go get golang.org/x/tools/cmd/godoc

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go get -d -v .

.PHONY: build $(SRCS)
build: deps
	CGO_ENABLED=0 GOOS=linux go build  -a -tags netgo -installsuffix netgo -ldflags $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install: deps
	go install -ldflags $(LDFLAGS)

.PHONY: lint
lint: $(LINT)
	@golint ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go test ./...

.PHONY: check
check: lint vet test build
