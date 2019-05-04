.PHONY: test
test:
	@go test ./...

.PHONY: ci-test
ci-test:
	@go test -v ./...
