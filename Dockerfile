FROM golang:1.13.7-buster

ARG app_env
ENV APP_ENV $app_env

# Enable go modules
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum into working directory
RUN go mod init workour-api

# Copy source code to image
COPY . .

# Download modules and "fresh" module
RUN go mod download

# PRODUCTION: build executable binary, set correct permissions,
# remove unneccesary files and execute binary
# DEVELOPMENT: download fresh package and run it
CMD if [ ${APP_ENV} = production ]; \
	then \
	export GIN_MODE=release && \
	cd .. && \
	go build ./app/ && \
	workour_app; \
	else \
	export GIN_MODE=debug && \
	go get github.com/pilu/fresh && \
	fresh; \
	fi

# Start with "fresh" on
EXPOSE 8080