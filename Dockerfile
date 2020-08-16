FROM golang:1.14.1-buster

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

# Download dependencies
RUN go mod download

# Install dockerize and wait for db to become available
RUN apt-get update && apt-get install -y wget

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Get package for hot reloading
RUN go get github.com/pilu/fresh

ENTRYPOINT dockerize -wait tcp://workour_db:5432 && fresh