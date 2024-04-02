.ONESHELL:

# Docker registry
SCW_ACCESS_KEY ?=
SCW_SECRET_KEY ?=

# Вычислаем версии
CURRENT_APP_VERSION := $(shell \
  git describe \
    --tags \
    --long \
    --always \
  | sed 's/-g.*$///'\
)

DOCKER_IMAGE_URL ?= rg.fr-par.scw.cloud/opdb/api:${CURRENT_APP_VERSION}

.PHONY: vars local_run.start_postgres local_run.stop_postgres tests.postgres.migrations
.PHONY: tests.postgres.test tests.postgres registry_login api.build api.push
.PHONY: terraform.plan terraform.apply

vars: ## Показать переменные
	: -------------------------------------------------------------------
	:  CURRENT_APP_VERSION: $(CURRENT_APP_VERSION)
	:  DOCKER_IMAGE_URL:    $(DOCKER_IMAGE_URL)
	: -------------------------------------------------------------------

local_run.start_postgres:
	docker run --name postgres-test -e POSTGRES_PASSWORD=qwerty -d -v ./sqls:/sqls postgres:15-alpine

local_run.stop_postgres:
	docker stop postgres-test
	docker rm -fv postgres-test

tests.postgres.migrations:
	: -------------------------------------------------------------------
	:  MIGRATIONS
	: -------------------------------------------------------------------
	docker exec --user postgres postgres-test psql -f /sqls/migrations/0000_init.sql

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
	make tests.postgres.test

	make local_run.stop_postgres

registry_login:
	docker login rg.fr-par.scw.cloud/opdb -u nologin -p ${SCW_SECRET_KEY}

api.build:
	docker build -t ${DOCKER_IMAGE_URL} --no-cache -f ./api/Dockerfile ./api

api.push:
	docker push ${DOCKER_IMAGE_URL}

terraform.plan:
	cd terraform
	terraform init
	terraform plan

terraform.apply:
	cd terraform
	terraform init
	terraform apply -auto-approve
