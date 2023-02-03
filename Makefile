## ---------- Variables ----------
BINARY_NAME=example
SERVICE_PORT?=3000
EXPORT_RESULT?=false 
LIST_SERVICES=catalog marketplace rating-service user-management
svc?=default

.PHONY: up down swag


## ---------- Docker ----------
build: ## Builds Docker services. Optional services. E.g: make build x y
	docker compose build $(filter-out $@,$(MAKECMDGOALS))

up: ## Builds and ups Docker services. Optional services. E.g: make up x y
	docker compose up -d --build $(filter-out $@,$(MAKECMDGOALS))
	watch docker compose ps

down: ## Brings Docker containers down. Requires services. E.g: make down x y
	docker compose rm -s -v $(filter-out $@,$(MAKECMDGOALS))


swag: ## Creates Swagger files. Requires service. E.g: make swag svc=x
ifeq ($(svc), $(filter $(svc), $(LIST_SERVICES))) 	
	if ! docker image inspect swagger-go >/dev/null 2>&1; then \
		echo "Container image for swagger does not exist. Building.."; \
		docker build . -t swagger-go -f doc/docker/dockerfile; \
	fi
	docker run --rm -v "$(shell pwd)/$(svc):/$(svc):ro" -v "$(shell pwd)/$(svc)/docs:/$(svc)/docs:rw" -w /$(svc) swagger-go swag init --parseInternal --parseDependency
else
	@echo "No service directory such as: $(svc)"
endif

swag-all: ## Creates Swagger files for all services
	$(MAKE) swag svc=catalog
	$(MAKE) swag svc=marketplace
	$(MAKE) swag svc=rating-service



# ---------- Shorcuts ----------
restart-full: ## Restarts one service while linting and updating swagger (down->lint->swagger->up). Requires service. E.g: make restart svc=x
ifeq ($(svc), $(filter $(svc), $(LIST_SERVICES)))
	$(MAKE) down $(svc)
	$(MAKE) lint-go svc=$(svc)
	$(MAKE) swag svc=$(svc)
	$(MAKE) up $(svc)
else
	@echo "No service directory such as: $(svc)"
endif


restart: ## Restarts services (down->up). Optional services. E.g make restart-all x y
	$(MAKE) down $(filter-out $@,$(MAKECMDGOALS))
	$(MAKE) up $(filter-out $@,$(MAKECMDGOALS))




## ---------- Linting ----------
lint-go: ## Use golintci-lint on your project. Requires service. E.g: make lint-go svc=x
ifeq ($(svc), $(filter $(svc), $(LIST_SERVICES)))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" = "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v "$(shell pwd)/$(svc):/app" -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)
else
	@echo "No service directory such as: $(svc)"
endif


lint-yaml: ## Use yamllint on the yaml file of your projects
ifeq ($(EXPORT_RESULT), true)
	go get -u github.com/thomaspoignant/yamllint-checkstyle
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	docker run --rm -it -v "$(shell pwd):/data" cytopia/yamllint -f parsable $(shell git ls-files '*.yml' '*.yaml') $(OUTPUT_OPTIONS)


help: ## Show Help Menu
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'