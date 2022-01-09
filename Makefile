SERVICE := metrika_datagen
CONTAINER_NAME := metrika_datagen_svc
VERSION := v0.0.1

HOST_PORT := 9000
JSON_FILE := ledger.json

.PHONY: build run logs stop

build:
	docker build . -t ${SERVICE}:${VERSION}

run:
	docker run --rm -d -p ${HOST_PORT}:8080 \
	--name ${CONTAINER_NAME} \
	-v $(shell pwd)/${JSON_FILE}:/app/${JSON_FILE} \
	 ${SERVICE}:${VERSION} -f ${JSON_FILE}

logs:
	docker logs -f ${CONTAINER_NAME}

stop:
	docker stop ${CONTAINER_NAME}