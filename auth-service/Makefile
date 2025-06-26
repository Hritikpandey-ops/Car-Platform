# Makefile for auth-service

APP_NAME=auth-service
DOCKER_COMPOSE=docker-compose
GO=go

# Load .env variables
include .env
export

.PHONY: all build run tidy fmt test docker-up docker-down docker-rebuild

all: build

## Build Go app
build:
	$(GO) build -o $(APP_NAME)

## Run app locally (non-docker)
run:
	$(GO) run main.go

## Tidy dependencies
tidy:
	$(GO) mod tidy

## Format Go code
fmt:
	$(GO) fmt ./...

## Run tests
test:
	$(GO) test ./...

## Start Docker services
docker-up:
	$(DOCKER_COMPOSE) up

## Stop Docker services
docker-down:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

## Rebuild Docker services
docker-rebuild:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) up --build


# Command	What it does
# make build	Compiles the Go binary
# make run	Runs the app locally
# make fmt	Formats Go code
# make tidy	Tidies up go.mod
# make test	Runs unit tests
# make docker-up	Starts services via Docker
# make docker-down	Stops services and removes volumes
# make docker-rebuild	Full rebuild of Docker containers