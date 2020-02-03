include .env

.EXPORT_ALL_VARIABLES:

.PHONY: create start stop

create-db:
	@docker create --name workour_db \
	-p 5432:${DATABASE_PORT} \
	-v ${PWD}/docker/data:/var/lib/postgresql/data \
	-e POSTGRES_USER=${DATABASE_USER} \
	-e POSTGRES_PASSWORD=${DATABASE_PSW} \
	-e POSTGRES_DB=${DATABASE_NAME} \
	postgres
	@docker start workour_db

build-app:
	@docker build -t workour_app \
	-f Dockerfile . \
	--build-arg app_env=development
	@docker run --name workour_app \
	-it -p 8080:8080 \
 	workour_app

start:
	@docker start workour_db

stop:
	@docker stop workour_db

test:
	go test ./tests