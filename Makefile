.PHONY: run.docker run.local deps clean test build.local build.linux.armv8 build.linux.armv7 build.linux build.osx build.windows

GOPKGS = $(shell go list ./... | grep -v /vendor/)

default: build.local

run.docker:
	docker build -t="leogregianin/brcep" .
	docker run -p 127.0.0.1:8000:8000/tcp leogregianin/brcep

run.local: deps
	go run $(GOPKGS) server.go

deps:
	go mod vendor

build.local:
	GO111MODULE=on go build -o bin/brcep .

build.linux.armv8:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GO111MODULE=on go build -o bin/brcep .

build.linux.armv7:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 GO111MODULE=on go build -o bin/brcep .

build.linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o bin/brcep .

build.osx:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o bin/brcep .

build.windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o bin/brcep .

test:
	go test -v $(GOPKGS) -coverprofile=coverage.txt -covermode=atomic

clean:
	rm -rf vendor
	rm -rf bin
