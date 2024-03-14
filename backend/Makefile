PROJECT := $(shell basename "$(PWD)")
GO_BIN := $(CURDIR)/build

GOPROXY ?= "https://proxy.golang.org,direct"

.PHONY: all registrar reporter clean

all: registrar reporter

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
clean:
	rm -rf ./build/*
