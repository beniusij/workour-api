include .env

.EXPORT_ALL_VARIABLES:

.PHONY: create start stop

build:
	@docker volume create --name=workour_db
	@docker-compose up --no-start

build-image:
	@docker build -t workour_api .

start:
	@docker-compose up

start-alone:
	@docker run -it -p 8080:8080 --rm -v ${PWD}/.:/app workour_api

clean:
	@docker image rm workour_api

stop:
	@docker-compose stop

test:
	go test ./tests

pre-deploy-test:
	go build -o bin/workour-api -v .
	heroku local