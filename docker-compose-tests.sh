#!/bin/bash
set -e
export PGVERSION=${PGVERSION:-16}
export DOCKER_COMPOSE_CMD=${DOCKER_COMPOSE_CMD:-"docker-compose"}
${DOCKER_COMPOSE_CMD} down --remove-orphans --rmi local || echo new or partial install
${DOCKER_COMPOSE_CMD} up -d postgres
${DOCKER_COMPOSE_CMD} up pgschema --no-deps --exit-code-from pgschema
${DOCKER_COMPOSE_CMD} up pgtester --no-deps --exit-code-from pgtester
