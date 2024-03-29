SHELL := /bin/bash

RUNDIR := $(shell pwd)

.PHONY: run
run: ### Start the docker-compose containers
	docker-compose --env-file ./config/.env build --force-rm 
	docker-compose --env-file ./config/.env  up --always-recreate-deps --remove-orphans 

.PHONY: test-retry
test-retry: ### Simulate a client request with retry 
	curl -X POST localhost:8001/reset_counters
	curl -s localhost:8001/stats | grep -e 'cluster.ingress_auth.upstream_rq_total' -e 'cluster.upstream.upstream_rq_total' -e 'cluster.upstream.upstream_rq_retry:' -e 'cluster.virtual_listener.upstream_rq_retry:'  -e 'cluster.virtual_listener.upstream_rq_total'
	curl -v  --header "x-envoy-retry-on: retriable-status-codes" --header "x-envoy-retriable-status-codes: 597" --header "x-auth-scenario-name: allow-and-add-header" --request GET http://localhost:8080/status/597
	curl -s localhost:8001/stats | grep -e 'cluster.ingress_auth.upstream_rq_total' -e 'cluster.upstream.upstream_rq_total' -e 'cluster.upstream.upstream_rq_retry:' -e 'cluster.virtual_listener.upstream_rq_retry:'  -e 'cluster.virtual_listener.upstream_rq_total'


# To start a network capture session, you will need first to run `make start-wireshark`
# that will listen on port 9999 for packets captures and start wireshark
# then run `make capture-envoy-traffic` to send packet captures to your local
# ip address on port 9999.
.PHONY: start-wireshark 
start-wireshark: ### Run wireshark listeneing for captured traffic on port 9999 
	nc -l -s 192.168.20.104 -p 9999 | sudo wireshark -k -S -i -

.PHONY: capture-envoy-traffic
capture-envoy-traffic: ### Capture envoy traffic and forward it to wireshark
	# Replace 192.168.20.104 by your local IP address
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-envoy-1) \
		--mount --uts --ipc --pid --net /bin/bash -c "/usr/sbin/tcpdump -n -i any -s 0 -w - not port 9999 | /bin/nc 192.168.20.104 9999"

.PHONY: help
help: ## Show this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
	@echo
