# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY:all smartraiden batchtransfer deploy newtestenv withdrawhash


export GOBIN = $(shell pwd)/build/bin
all: smartraiden batchtransfer deploy newtestenv  withdrawhash
smartraiden:
	go  install ./cmd/smartraiden
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





