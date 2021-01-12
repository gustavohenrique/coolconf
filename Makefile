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
	$(GO) get -u golang.org/x/tools/cmd/goimports
