SHELL := /bin/bash

RUNDIR := $(shell pwd)

run: 
	docker-compose --env-file ./config/.env  up --build --force-recreate --no-deps 

test-mismatch:
	http -v localhost:8080/headers x-auth-scenario-name:transition-deny
	@echo "Checking envoy logs... grep for log..."
	docker logs envoy 2>&1 | grep 'script log' 

test-match:
	http -v localhost:8080/headers x-auth-scenario-name:transition-allow
	@echo "Checking envoy logs... grep for log..."
	docker logs envoy 2>&1 | grep 'script log' 

test-no-extauthz-status:
	http -v localhost:8080/headers x-auth-scenario-name:allow
	@echo "Checking envoy logs... grep for log..."
	docker logs envoy 2>&1 | grep 'script log' 
