SHELL := /bin/bash

RUNDIR := $(shell pwd)

# You will need to run `make cert` one time so that the go-httpbin service
# inside the tunnel container works as a https server
.PHONY: cert
cert: ### Make certificate to use by the server container (go-httpbin)
	openssl genrsa -out /tmp/server.key 2048
	openssl ecparam -genkey -name secp384r1 -out /tmp/server.key
	openssl req -new -x509 -sha256 -key /tmp/server.key -out /tmp/server.crt -days 3650

.PHONY: run
run: ### Start the docker-compose containers
	docker-compose build --force-rm 
	docker-compose --env-file ./config/.env  up --always-recreate-deps --remove-orphans 

.PHONY: test-tunnel
test-tunnel: ### Simulate a CONNECT request via HTTP proxy but expected to go to the tunnel service
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-client-1) \
		--mount --uts --ipc --pid --net \
		/bin/sh -c "SSLKEYLOGFILE=/tmp/tlskey curl -k -v -x envoy:10000 https://foo.tunnel.example.com/headers"

.PHONY: test-outbound
test-outbound: ### Simulate a CONNECT request via HTTP proxy but expected to go to the outbound proxy
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-client-1) \
		--mount --uts --ipc --pid --net \
		/bin/sh -c "SSLKEYLOGFILE=/tmp/tlskey curl -k -v -x envoy:10000 https://www.google.com"

# To start a network capture session, you will need first to run `make start-wireshark`
# that will listen on port 9999 for packets captures and start wireshark
# then run `make capture-envoy-traffic` to send packet captures to your local
# ip address on port 9999.
.PHONY: start-wireshark 
start-wireshark: ### Run wireshark listeneing for captured traffic on port 9999 
	nc -l -p 9999 | sudo wireshark -k -S -i -

.PHONY: capture-envoy-traffic
capture-envoy-traffic: ### Capture envoy traffic and forward it to wireshark
	# Replace 192.168.20.104 by your local IP address
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-envoy-1) \
		--mount --uts --ipc --pid --net \
		/bin/bash -c "/usr/sbin/tcpdump -n -i any -s 0 -w - not port 9999 | /bin/nc 192.168.20.104 9999"

.PHONY: help
help: ## Show this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
	@echo
