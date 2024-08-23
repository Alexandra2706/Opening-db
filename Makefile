.ONESHELL:

# Docker registry
SCW_ACCESS_KEY ?=
SCW_SECRET_KEY ?=

# Вычислаем версии
CURRENT_APP_VERSION ?= $(shell \
  git describe \
    --tags \
    --long \
    --always \
  | sed 's/-g.*$///'\
)

DOCKER_IMAGE_URL ?= rg.fr-par.scw.cloud/opdb/api:${CURRENT_APP_VERSION}
DOCKER_IMAGE_URL_SYNC ?= rg.fr-par.scw.cloud/opdb/sync:${CURRENT_APP_VERSION}

.PHONY: vars local_run.start_postgres local_run.stop_postgres tests.postgres.migrations
.PHONY: tests.postgres.test tests.postgres registry_login api.build api.push
.PHONY: tests.sync.run tests.sync sync.build
.PHONY: terraform.plan terraform.apply

vars: ## Показать переменные
	: -------------------------------------------------------------------
	:  CURRENT_APP_VERSION: $(CURRENT_APP_VERSION)
	:  DOCKER_IMAGE_URL:    $(DOCKER_IMAGE_URL)
	: -------------------------------------------------------------------

local_run.start_postgres:
	docker run --name postgres-test -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d -v ./sqls:/sqls postgres:15-alpine

local_run.stop_postgres:
	docker stop postgres-test
	docker rm -fv postgres-test

tests.postgres.migrations:
	: -------------------------------------------------------------------
	:  MIGRATIONS
	: -------------------------------------------------------------------
	docker exec --user postgres postgres-test psql -f /sqls/migrations/0000_init.sql

tests.postgres.data_create:
	: -------------------------------------------------------------------
	:  DATA CREATE
	: -------------------------------------------------------------------
	docker exec --user postgres postgres-test psql -f /sqls/datagenerator/generator.sql

tests.postgres.test:
	: -------------------------------------------------------------------
	:  TEST: table1.sql
	: -------------------------------------------------------------------
	#docker exec --user postgres postgres-test psql -f /sqls/tests/table1.sql

tests.postgres:
	make local_run.start_postgres
	sleep 10
	make tests.postgres.migrations
	make tests.postgres.migrations # Повторно применяем миграции для проверки IF EXISTS
	make tests.postgres.data_create
	make tests.postgres.test

	make local_run.stop_postgres

tests.sync.run:
	make sync.build
	docker run --name sync -it --net=host -e DB_URL=postgres://postgres:qwerty@localhost:5432/postgres ${DOCKER_IMAGE_URL_SYNC}
	docker rm sync

tests.sync:
	make local_run.start_postgres
	sleep 10
	make tests.postgres.migrations

	make tests.sync.run
	make local_run.stop_postgres

tests.unit:
	cd api
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	go generate
	go test ./...

registry_login:
	docker login rg.fr-par.scw.cloud/opdb -u nologin -p ${SCW_SECRET_KEY}

api.build:
	docker build -t ${DOCKER_IMAGE_URL} --no-cache -f ./api/Dockerfile ./api

sync.build:
	docker build -t ${DOCKER_IMAGE_URL_SYNC} --no-cache -f ./jobs/shikimori-sync/Dockerfile ./jobs/shikimori-sync

api.push:
	docker push ${DOCKER_IMAGE_URL}

terraform.plan:
	cd terraform
	TF_VAR_api_container=${DOCKER_IMAGE_URL} terraform init
	TF_VAR_api_container=${DOCKER_IMAGE_URL} terraform plan

terraform.apply:
	cd terraform
	TF_VAR_api_container=${DOCKER_IMAGE_URL} terraform init
	TF_VAR_api_container=${DOCKER_IMAGE_URL} terraform apply -auto-approve
