SHELL := /bin/bash
.DEFAULT_GOAL := help
COMPOSE ?= docker compose

## Start all services
up:
	$(COMPOSE) up -d

## Stop all services
down:
	$(COMPOSE) down

## One-time bootstrap: buckets, topics, defaults, demo data
init:
	bash scripts/init.sh

## Submit Spark job(s)
spark-submit:
	bash scripts/spark_submit.sh

## Great Expectations validations
ge:
	bash scripts/run_ge.sh

## dbt models + docs
dbt:
	bash scripts/run_dbt.sh

## Tail logs
logs:
	$(COMPOSE) logs -f

## Status
ps:
	$(COMPOSE) ps

## Nuke everything (volumes, orphans)
clean:
	$(COMPOSE) down -v --remove-orphans
	docker volume prune -f
	docker network prune -f

## Help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS=":.*?## "}; {printf "\033[36m%-14s\033[0m %s\n", $$1, $$2}'
