install-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(CURDIR) v1.61.0

lint:
	$(CURDIR)/golangci-lint run --verbose

test:
	@echo "Running tests with coverage..."
	@go test -race -cover ./...