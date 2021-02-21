#!/usr/bin/make -f

all: clean test build install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

build_test_container:
	docker-compose -f ./deploy/test/docker-compose.yml --project-directory . build

start_test_containers:
	docker-compose -f ./deploy/test/docker-compose.yml --project-directory . up

stop_test_containers:
	docker-compose -f ./deploy/test/docker-compose.yml --project-directory . down

clean:
	rm -f ebrelayer

install:
	go install ./cmd/ebrelayer

.PHONY: all build go-mod-cache build_test_container start_test_containers stop_test_containers clean install test lint all
