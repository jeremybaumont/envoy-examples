# External authorization HTTP filter dry run (lua log)

## Overview

This is a sandbox with envoy as front proxy and [httpbin](https://eu.httpbin.org/) as upstream service.
We enable a external authorization HTTP filter in the HTTP filter chain to a gRPC service.
We use the external authorization HTTP filter in dry run mode. Meaning, it allows the requests but add a x-extauthz-status-code header with the status it would have set in wet mode.
We use the `x-auth-scenario-name` to control the value of the `x-extauthz-status-code` header.
We enable HTTP lua filter to log when `x-extauthz-status-code` differs from the upstream response code.

![demo screencast](./wasm-mismatch.gif "demo")

## Requirement

* Install make
* Install docker
* Install docker-compose

## Usage

To run envoy, the external authorization gRPC service and upstream service: 
```
make run
```

## Test

### Test mismatch 
This test checks when there is a mismatch between the status code that would have been returned by the ext_authz HTTP filter, and the actual upstream response code.
In this scenario, the x-auth-status-code header is appended to the request by the ext authz HTTP filter with the value 403. 
The upstream response status code is 200.  

So there is a mismatch. A log is written by the HTTP lua filter 

The makefile has the following rule to run the mismatch test:
```
test-mismatch:
	http -v localhost:8080/headers x-auth-scenario-name:transition-deny
	@echo "Dumping envoy stats... grep for authMismatch metric name..."
	http localhost:8001/stats | grep authMismatch
```

To trigger the test, run the following command:
```
make test-mismatch
```

### Test match

This test checks when there is a match between the status code that would have been returned by the ext_authz HTTP filter, and the actual upstream response code.
In this scenario, the x-auth-status-code header is appended to the request by the ext authz HTTP filter with the value 200. 
The upstream response status code is 200.  
So there is a match. No log is written. 

The makefile has the following rule to run the mismatch test:
```
test-match:
	http -v localhost:8080/headers x-auth-scenario-name:transition-allow
	@echo "Dumping envoy stats... grep for authMismatch metric name..."
	http localhost:8001/stats | grep authMismatch
```

To trigger the test, run the following command:
```
make test-match
```


### Test when no x-auth-status-code exists

This test checks when there is no status code is appended by the ext_authz HTTP filter.
In this scenario, the x-auth-status-code header is not appended to the request by the ext authz HTTP.
The upstream response status code is 200.  
No log is written. 


```
test-no-extauthz-status:
	http -v localhost:8080/headers x-auth-scenario-name:allow
	@echo "Dumping envoy stats... grep for authMismatch metric name..."
	http localhost:8001/stats | grep authMismatch
```                                                                                                                                                                                                      

To trigger the test, run the following command:
```
make test-no-extauthz-status
```


