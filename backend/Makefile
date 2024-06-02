PROJECT := $(shell basename "$(PWD)")
GO_BIN := $(CURDIR)/build

GOPROXY ?= "https://proxy.golang.org,direct"

.PHONY: all registrar reporter generator clean

all: registrar reporter generator

registrar:
	GOPROXY=$(GOPROXY) \
	go build -ldflags "-s -w" \
	-o $(GO_BIN)/registrar \
	./cmd/registrar

reporter:
	GOPROXY=$(GOPROXY) \
	go build -ldflags "-s -w" \
	-o $(GO_BIN)/reporter \
	./cmd/reporter

generator:
	GOPROXY=$(GOPROXY) \
	go build -ldflags "-s -w" \
	-o $(GO_BIN)/gen \
	./cmd/generator

clean:
	rm -rf ./build/*
