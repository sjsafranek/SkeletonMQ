##=======================================================================##
## Makefile
## Created: Wed Aug 05 14:35:14 PDT 2015 @941 /Internet Time/
# :mode=makefile:tabSize=3:indentSize=3:
## Purpose:
##======================================================================##

SHELL=/bin/bash
PROJECT_NAME = GeoSkeletonServer
GPATH = $(shell pwd)

.PHONY: fmt get-deps update-deps test install build scrape clean

install: fmt get-deps
	@GOPATH=${GPATH} go build -o server *.go

build: fmt get-deps
	@GOPATH=${GPATH} go build -o server *.go

get-deps:
	mkdir -p "src"
	mkdir -p "pkg"
	mkdir -p "log"
	@GOPATH=${GPATH} go get github.com/cihub/seelog
	@GOPATH=${GPATH} go get github.com/gorilla/mux
	@GOPATH=${GPATH} go get github.com/gorilla/websocket
	@GOPATH=${GPATH} go get github.com/schollz/jsonstore

update-deps: get-deps
	@GOPATH=${GPATH} go get -u github.com/cihub/seelog
	@GOPATH=${GPATH} go get -u github.com/gorilla/mux
	@GOPATH=${GPATH} go get -u github.com/gorilla/websocket
	@GOPATH=${GPATH} go get -u -v github.com/schollz/jsonstore


fmt:
	@GOPATH=${GPATH} gofmt -s -w *.go

test:
	##./tcp_test.sh
	./benchmark.sh

scrape:
	@find src -type d -name '.hg' -or -type d -name '.git' | xargs rm -rf

clean:
	@GOPATH=${GPATH} go clean
