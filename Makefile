ifneq ("$(wildcard .env.local)", "")
	include .env.local
	export
endif

.PHONY: test_all 
test_all: clean_test lint test integration

.PHONY: lint
lint:
	golangci-lint run --verbose

.PHONY: test
test:
	go test ./... -skip=TestIntegration -test.v

.PHONY: integration
integration:
	@test -n "$(API_KEY)" || (echo 'API_KEY env required to run integration tests' && exit 1)
	go test . -run=TestIntegration -test.v

.PHONY: clean_test
clean_test:
	go clean -testcache
