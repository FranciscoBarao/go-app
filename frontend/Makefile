## ---------- Variables ----------
BUILD_DIR := dist
COVERAGE_DIR := coverage

# Installs dependencies when the lockfile changes. Other targets depend on it.
node_modules: package-lock.json
	npm ci

## ---------- Linting ----------
lint: node_modules ## Runs Biome checks
	npm run lint


## ---------- Build ----------
build: node_modules ## Typechecks and builds into dist/
	npm run build

build-docker: ## Builds the production nginx image
	docker compose build frontend


## ---------- Dev ----------
dev: node_modules ## Starts the Vite dev server on the host
	npm run dev


## ---------- Test ----------
test: node_modules ## Runs tests once
	npm run test

test/cover: node_modules ## Runs tests with coverage into coverage/
	npm run test:cover


## ---------- Shortcuts ----------
clean: ## Removes build and coverage output
	rm -rf $(BUILD_DIR)/ $(COVERAGE_DIR)/

check: clean build test lint ## Runs build, tests and lint

run: build-docker ## Rebuilds and restarts the frontend container
	docker compose down frontend
	docker compose up -d frontend


help: ## Show Help Menu
	@grep -E '^[a-zA-Z._/-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
