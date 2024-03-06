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

DOCKER_IMAGE_NAME ?= rg.fr-par.scw.cloud/opdb/api:${CURRENT_APP_VERSION}

vars: ## Показать переменные
	: -------------------------------------------------------------------
	:  CURRENT_APP_VERSION: $(CURRENT_APP_VERSION)
	:  DOCKER_IMAGE_URL:    $(DOCKER_IMAGE_URL)
	: -------------------------------------------------------------------

registry_login:
	docker login rg.fr-par.scw.cloud/opdb -u nologin -p ${SCW_SECRET_KEY}

api.build:
	docker build -t ${DOCKER_IMAGE_URL} --no-cache -f ./api/Dockerfile ./api

api.push:
	docker push ${DOCKER_IMAGE_URL}
