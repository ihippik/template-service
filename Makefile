.PHONY: test lint docs help

$(GOLANGCI):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@

$(SWAG):
	go install github.com/swaggo/swag/cmd/swag

docs: $(SWAG) ## generate swag docs
	swag init --parseVendor -g main.go

test: ## run unit (short) tests
	go test -short ./...

lint: $(GOLANGCI) ## Check the project with lint
	golangci-lint run

help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

