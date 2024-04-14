ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=postgres password=top_secret dbname=gohw_test host=localhost port=5432 sslmode=disable
endif

ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=top_secret dbname=gohw host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal
MOCKGEN_TAG=1.6.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations


.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down


.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: init-pg
init-pg:
	docker run --env-file ./.local.env --mount type=volume,target=/var/lib/postgresql/data -p 127.0.0.1:5433:5432/tcp -d postgres:16-alpine

.PHONY: run
run:
	env $$(cat .local.env | xargs) go run cmd/main.go http

.PHONY: kill-9000
kill:
	kill -9 $$(lsof -ti :9000)

.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_TAG)
endif

.PHONY: .generate-mockgen
.generate-mockgen:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...

.PHONY: test-module
test-module:
	env $$(cat .local.env | xargs) go test ./internal/...

.PHONY: test-integration
test-integration:
	env $$(cat .local.env | xargs) go test -tags=integration ./tests/...

# Kafka
build:
	docker-compose build

up-all:
	docker-compose up -d zookeeper kafka1 kafka2 kafka3

down:
	docker-compose down

run:
	go run ./cmd/route-kafka



