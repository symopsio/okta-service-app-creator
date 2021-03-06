# Goals:
# - user can build binaries on their system without having to install special tools
# - user can fork the canonical repo and expect to be able to run CircleCI checks
#
# This makefile is meant for humans

VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.Version=$(VERSION)"'

#test:
#	GO111MODULE=on go test -mod=vendor -v ./...

all: dist/okta-service-app-creator-$(VERSION)-darwin-amd64 dist/okta-service-app-creator-$(VERSION)-linux-amd64

clean:
	rm -rf ./dist

dist/:
	mkdir -p dist

dist/okta-service-app-creator-$(VERSION)-darwin-amd64: | dist/
	(cd cmd && GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -mod=vendor $(LDFLAGS) -o ../$@)

dist/okta-service-app-creator-$(VERSION)-linux-amd64: | dist/
	(cd cmd && GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor $(LDFLAGS) -o ../$@)

.PHONY: clean all
