# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY:all smartraiden batchtransfer deploy newtestenv withdrawhash


export GOBIN = $(shell pwd)/build/bin
export GIT_COMMIT=$(shell git rev-list -1 HEAD)
export GO_VERSION=$(shell go version|sed 's/ //g')
export BUILD_DATE=$(shell date|sed 's/ //g')
all: smartraiden batchtransfer deploy newtestenv  withdrawhash
smartraiden:
	go  install -ldflags ' -X github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl.GitCommit=$(GIT_COMMIT) -X github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl.GoVersion="$(GO_VERSION)" -X github.com/SmartMeshFoundation/SmartRaiden/cmd/smartraiden/mainimpl.BuildDate="$(BUILD_DATE)" '  ./cmd/smartraiden
	@echo "Done building."
	@echo "Run \"$(GOBIN)/smartraiden\" to launch smartraiden."

batchtransfer:
	go  install ./cmd/tools/batchtransfer
	@echo "Done building."
	@echo "Run \"$(GOBIN)/batchtransfer\" to launch batchtransfer."

deploy:
	go  install ./cmd/tools/deploy
	@echo "Done building."
	@echo "Run \"$(GOBIN)/deploy\" to launch deploy."

newtestenv:
	go  install ./cmd/tools/newtestenv
	@echo "Done building."
	@echo "Run \"$(GOBIN)/newtestenv\" to launch newtestenv."

withdrawhash:
	go  install ./cmd/tools/withdrawhash
	@echo "Done building."
	@echo "Run \"$(GOBIN)/withdrawhash\" to launch withdrawhash."





