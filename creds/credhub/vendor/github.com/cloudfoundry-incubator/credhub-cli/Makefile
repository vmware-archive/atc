.PHONY: all build ci clean dependencies format ginkgo test

ifeq ($(GOOS),windows)
DEST = build/credhub.exe
else
DEST = build/credhub
endif

ifndef VERSION
VERSION = dev
endif

GOFLAGS := -o $(DEST) -ldflags "-X github.com/cloudfoundry-incubator/credhub-cli/version.Version=${VERSION}"

all: dependencies test clean build

clean:
	rm -rf build

dependencies:
	go get github.com/onsi/ginkgo/ginkgo
	go get golang.org/x/tools/cmd/goimports
	go get github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/kardianos/govendor
	govendor sync

format:
	goimports -w .
	go fmt .

ginkgo:
	ginkgo -r -randomizeSuites -randomizeAllSpecs -race -p 2>&1

test: format ginkgo

ci: dependencies ginkgo

build:
	mkdir -p build
	go build $(GOFLAGS)
