# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY:all Photon batchtransfer deploy newtestenv withdrawhash


export GOBIN = $(shell pwd)/build/bin
export GIT_COMMIT=$(shell git rev-list -1 HEAD)
export GO_VERSION=$(shell go version|sed 's/ //g')
export BUILD_DATE=$(shell date|sed 's/ //g')
export VERSION=0.91
all: Photon batchtransfer deploy newtestenv  withdrawhash
Photon:
	go  install -ldflags ' -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GitCommit=$(GIT_COMMIT) -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.GoVersion="$(GO_VERSION)" -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.BuildDate="$(BUILD_DATE)" -X github.com/SmartMeshFoundation/Photon/cmd/photon/mainimpl.Version="$(VERSION)" '  ./cmd/photon
	@echo "Done building."
	@echo "Run \"$(GOBIN)/photon\" to launch photon."

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





