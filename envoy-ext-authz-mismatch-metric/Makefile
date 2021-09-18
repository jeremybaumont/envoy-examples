SHELL := /bin/bash

# check wether compiled with wasm-pack
WASMDIR := ./wasm-extauthz-mismatch
WASM_PATH := $(WASMDIR)/pkg/extauthz_mismatch_bg.wasm

RUNDIR := $(shell pwd)

# a small optimized binary without debug info, useful for releases
build-wasm: clean-wasm
	wasm-pack build --release $(WASMDIR)

run: build-wasm
	WASM_PATH=$(WASM_PATH) docker-compose --env-file ./config/.env  up --build --force-recreate --no-deps 

clean-wasm:
	cargo clean --target-dir $(WASMDIR)/target --manifest-path $(WASMDIR)/Cargo.toml
	rm -rf $(WASMDIR)/pkg

test-mismatch:
	http -v localhost:8080/headers x-auth-scenario-name:transition-deny
	@echo "Dumping envoy stats... grep for authMismatch metric name..."
	http localhost:8001/stats | grep authMismatch

test-match:
	http -v localhost:8080/headers x-auth-scenario-name:transition-allow
	@echo "Dumping envoy stats... grep for authMismatch metric name..."
	http localhost:8001/stats | grep authMismatch

test-no-extauthz-status:
	http -v localhost:8080/headers x-auth-scenario-name:allow
	@echo "Dumping envoy stats... grep for authMismatch metric name..."
	http localhost:8001/stats | grep authMismatch
