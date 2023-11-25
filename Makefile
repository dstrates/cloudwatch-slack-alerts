.PHONY: build clean

# set the version as the latest commit sha if it's not already defined
ifndef VERSION
# check if there are code changes that aren't commited
# add a -tainted label to the end of the version if there are
ifneq ($(shell git status --porcelain), )
TAINT := -tainted
endif
VERSION := $(shell git rev-list -1 HEAD)$(TAINT)
endif

GOENV := CGO_ENABLED=0
GOFLAGS := -ldflags "-X 'github.com/nullify-platform/logger/pkg/logger.Version=$(VERSION)'"

build:
	$(GOENV) go build $(GOFLAGS) -o bin/main ./cmd/...

package:
	$(GOENV) GOOS=linux GOARCH=amd64 go build -o bin/main ./cmd/...
	zip -j bin/alerter.zip bin/main

clean:
	rm -rf ./bin ./vendor Gopkg.lock coverage.*

tidy:
	go mod tidy

format:
	gofmt -w ./...

lint:
	docker build --quiet --target golangci-lint -t golangci-lint:latest .
	docker run --rm -v $(shell pwd):/app -w /app golangci-lint golangci-lint run ./...

unit:
	go test -v -skip TestIntegration ./...
