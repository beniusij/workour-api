include .env

.EXPORT_ALL_VARIABLES:

.PHONY: create start stop

build:
	@docker volume create --name=workour_db
	@docker-compose up --no-start

start:
	@docker-compose up

stop:
	@docker-compose stop

test:
	go test ./tests