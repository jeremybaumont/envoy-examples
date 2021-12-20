SHELL := /bin/bash

RUNDIR := $(shell pwd)

run: 
	docker-compose --env-file ./config/.env  up --build --force-recreate --no-deps 
