## ---------- Variables ----------
BUILD_DIR := build
GOLANGCI_LINT_VERSION ?= v2.11.4


## ---------- Linting ----------
lint: ## Runs golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	$$(go env GOPATH)/bin/golangci-lint run --concurrency 5 -v ./...

## ---------- Build ----------
build:
	go build -o=$(BUILD_DIR)/catalog main.go

build-docker:
	docker compose build catalog catalog-db

clean-build:
	rm -rf $(BUILD_DIR)/

## ---------- Mock ----------
build-mock:
	go install go.uber.org/mock/mockgen@latest
	go generate -x ./...

clean-mock:
	find . -type f -name '*_mock.go' -exec rm {} +


## ---------- Test ----------
test: ## Runs tests
	go test -race -buildvcs ./...

test/cover: ## Runs tests with coverage
	go test -v -race -buildvcs -coverprofile=$(BUILD_DIR)/coverage.out ./...


## ---------- Others ----------
tidy: ## Runs fmt & mod tidy
	go fmt ./...
	go mod tidy -v


## ---------- Shortcuts ----------

clean: clean-mock clean-build


check: clean build-mock tidy test lint

run: build-docker ## Up/Down with build
	docker compose down catalog catalog-db
	docker compose up -d catalog catalog-db
    ## watch docker compose ps


help: ## Show Help Menu
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
