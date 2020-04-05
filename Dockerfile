FROM golang:1.14.1-buster

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

# Download dependencies
RUN go mod download

# Get package for hot reloading
RUN go get github.com/pilu/fresh

ENTRYPOINT fresh