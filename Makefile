.PHONY: build
build: ## Builds the executable
	go build -o generate .

.PHONY: test
test: ## Runs tests
	go test -v -p 1 ./...

.PHONY: quality
quality: ## Runs quality checks (golangci-lint, govulncheck, and tests)
	@echo "--> Formatting code..."
	go fmt ./...
	@echo "--> Running golangci-lint..."
	golangci-lint run ./...
	@echo "--> Running go vet.."
	go vet ./...
	@echo "--> Checking vulnerabilities.."
	govulncheck ./...
	@echo "--> Running tests..."
	$(MAKE) test


.PHONY: clean
clean: ## Cleans temporary files
	rm -rf generate cover.out dist

.PHONY: help
help:	## Show this help
	@grep -hE '^[A-Za-z0-9_ \-]*?:.*##.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

