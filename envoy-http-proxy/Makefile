SHELL := /bin/bash

RUNDIR := $(shell pwd)
LOCAL_IP_ADDRESS := $(shell ip addr show wlan0 | grep inet | awk '{print $$2}' | awk -F'/' '{print $$1}')  

.PHONY: run
run: ### Start the docker-compose containers
	docker-compose --env-file ./config/.env build --force-rm 
	docker-compose --env-file ./config/.env  up --always-recreate-deps --remove-orphans 

.PHONY: test-tunnel-https
test-tunnel-https: ### Simulate a CONNECT request via HTTP proxy but expected to go to the tunnel service
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-client-1) \
		--mount --uts --ipc --pid --net \
		/bin/sh -c "SSLKEYLOGFILE=/tmp/tlskey curl -k -v -x envoy:10000 https://foo.tunnel.example.com/headers"

.PHONY: test-outbound-https
test-outbound-https: ### Simulate a CONNECT request via HTTP proxy but expected to go to the outbound proxy
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-client-1) \
		--mount --uts --ipc --pid --net \
		/bin/sh -c "SSLKEYLOGFILE=/tmp/tlskey curl -k -v -x envoy:10000 https://www.google.com"

.PHONY: test-tunnel-http
test-tunnel-http: ### Simulate a CONNECT request via HTTP proxy but expected to go to the tunnel service
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-client-1) \
		--mount --uts --ipc --pid --net \
		/bin/sh -c "SSLKEYLOGFILE=/tmp/tlskey curl -k -v -x envoy:10000 http://foo.tunnel.example.com/headers"

.PHONY: test-outbound-http
test-outbound-http: ### Simulate a CONNECT request via HTTP proxy but expected to go to the outbound proxy
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-client-1) \
		--mount --uts --ipc --pid --net \
		/bin/sh -c "SSLKEYLOGFILE=/tmp/tlskey curl -k -v -x envoy:10000 http://www.google.com"

# To start a network capture session, you will need first to run `make start-wireshark`
# that will listen on port 9999 for packets captures and start wireshark
# then run `make capture-envoy-traffic` to send packet captures to your local
# ip address on port 9999.
.PHONY: start-wireshark 
start-wireshark: ### Run wireshark listening for captured traffic on port 9999 
	nc -l -s $(LOCAL_IP_ADDRESS) -p 9999 | sudo wireshark -k -S -i -

.PHONY: capture-envoy-traffic
capture-envoy-traffic: ### Capture envoy traffic and forward it to wireshark
	sudo nsenter -t $$(docker inspect -f {{.State.Pid}} envoy-http-proxy-envoy-1) \
		--mount --uts --ipc --pid --net /bin/bash -c "/usr/sbin/tcpdump -n -i any -s 0 -w - not port 9999 | /bin/nc $(LOCAL_IP_ADDRESS) 9999"

.PHONY: help
help: ## Show this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
	@echo
