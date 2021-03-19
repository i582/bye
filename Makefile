GOPATH_DIR=`go env GOPATH`

check:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.30.0
	@echo "running linters..."
	@$(GOPATH_DIR)/bin/golangci-lint run ./...
#	@echo "running tests..."
#	@go test -tags tracing -count 3 -race -v ./src/...
#	@go test -race ./example/custom
	@echo "everything is OK"

.PHONY: check
