.DEFAULT_GOAL := build
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PACKAGES = $(shell find ./ -type d -not -path "./vendor" -not -path "./vendor/*" -not -path "./.git" -not -path "./.git/*" -not -path "./config/files")
SHELL=/bin/bash

clean:
	@rm -f ion-cli
	@rm *.out

test:
	@go test ./... -v

format:
	@gofmt -s -w ${GOFILES_NOVENDOR}

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

build:
	@go get -t -v ./...
	@go build -o ion-cli .

check:
	@if [ -n "$(shell gofmt -l ${GOFILES_NOVENDOR})" ]; then \
		echo 1>&2 'The following files need to be formatted:'; \
		gofmt -l .; \
		exit 1; \
		fi

vet:
	@go vet ${GOFILES_NOVENDOR}

lint:
	@golint ${GOFILES_NOVENDOR}
