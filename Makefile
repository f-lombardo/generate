.PHONY: build
build: ## Builds the executable
	go build -o generate .

.PHONY: test
test: ## Runs tests
	go test -v -p 1 ./...

.PHONY: clean
clean: ## Cleans temporary files
	rm -rf generate cover.out dist

.PHONY: help
help:	## Show this help
	@grep -hE '^[A-Za-z0-9_ \-]*?:.*##.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

