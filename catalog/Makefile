## ---------- Variables ----------
EXPORT_RESULT?=true 
BUILD_DIR := build


## ---------- Linting ----------
lint: ## Uses golintci-lint
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" = "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v "$(shell pwd):/app" -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)


## ---------- Build ----------
build-exe: ## Builds exe
	go build -o=$(BUILD_DIR)/catalog main.go

build-docker: ## Builds Docker services
	docker compose build catalog catalog-db

clean-build: ## Deletes exe build
	rm -rf $(BUILD_DIR)/


## ---------- Mock ----------
mock: ## Generates mocks
	go get github.com/golang/mock/mockgen@latest
	go install github.com/golang/mock/mockgen
	go generate -x ./...


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
check: mock clean-build tidy test lint ## Verifies everything

run: build-docker ## Up/Down with build
	docker compose down catalog catalog-db
	docker compose up -d catalog catalog-db
    ## watch docker compose ps


help: ## Show Help Menu
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'