.PHONY: test

GO := $(shell which go)
GOTEST := $(shell which gotest)

test: tests
tests:
	$(GOTEST) -v -failfast -cover ./...

lint: linter
linter:
	goimports -w .

install: setup
setup:
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) get -u github.com/rakyll/gotest
	$(GO) install github.com/rakyll/gotest
