include .env

.EXPORT_ALL_VARIABLES:

.PHONY: create-db

create-db:
	@docker create --name workour_db \
	-p 5432:${DATABASE_PORT} \
	-v ${PWD}/docker/data:/var/lib/postgresql/data \
	-e POSTGRES_USER=${DATABASE_USER} \
	-e POSTGRES_PASSWORD=${DATABASE_PSW} \
	-e POSTGRES_DB=${DATABASE_NAME} \
	postgres
	@docker start workour_db