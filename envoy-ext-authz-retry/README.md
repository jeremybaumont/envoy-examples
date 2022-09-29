# External authorization HTTP filter dry run 

## Overview

This is a sandbox with envoy as front proxy and [httpbin](https://eu.httpbin.org/) as upstream service.
We enable a external authorization HTTP filter in the HTTP filter chain to a gRPC service.
We use the external authorization HTTP filter in dry run mode. Meaning, it allows the requests but add a x-extauthz-status-code header with the status it would have set in wet mode.
We use the `x-auth-scenario-name` to control the value of the `x-extauthz-status-code` header.
We enable also a wasm HTTP filter that will store the value of the `x-extauthz-status-code` header, and will compare it to the upstream response status code.
If they mismatches, it will emits a metric.

![demo screencast](./wasm-mismatch.gif "demo")

## Requirement

* Install make
* Install docker
* Install docker-compose
* Install wasm-pack

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
So there is a mismatch. A metric is emitted, authMistmach counter metric is incremented.

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

The following diagram sums up the interactions:
                                                                                                                                                                                                         
                                                        Envoy                                                                                                                                             
                                                                                                                                                                                                          
                                                                                                                                                                                                          
                                    +--------------------------------------------------+                                                                                                                  
                                    |                                                  |                                                                                                                  
                                    |                                                  |                                                                                                                  
                                    |  ext_authz filter         wasm mismatch filter   |                       httpbin                                                                                    
                                    | +-----------------+        +----------------+    |                    +-----------+                                                                                 
                  1. test mismatch  | |                 |        |4. store        |    |                    |           |                                                                                 
                            ---------->                 -------->|ext_authz_status------------------------->|           |                                                                                 
                                    | |                 |        |403             |    |                    |           |                                                                                 
                                    | |                 |        |                |    |                    |           |                                                                                 
                           <-----------                 <--------|                |    |                    |           |                                                                                 
                                    | |                 |        |                <-------------------------|           |                                                                                 
                                    | |                 |        |                |    |                    +-----------+                                                                                 
                                    | +----|------------+        +----------------+    |  5. response status code 200                                                                                     
                                    |      |    ^                  6. ext_authz_status |                                                                                                                  
                                    |      |    |                  != response status  |                                                                                                                  
                                    +-----------|-----------------------------------|--+                                                                                                                  
                                   2. Check|    |3. allow w/                        |                                                                                                                     
                                           |    |x-auth-status-code:403        7. emit authMismatch counter                                                                                               
                                      +----v----|--------+                          |                                                                                                                     
                                      |                  |                          |                                                                                                                     
                                      |                  |                          |                                                                                                                     
                                      |                  |                          |                                                                                                                     
                                      +------------------+                          v                                                                                                                     
                                                                                                                                                                                                          
                                  ext authz gRPC service                                                                                                                                                  

### Test match

This test checks when there is a match between the status code that would have been returned by the ext_authz HTTP filter, and the actual upstream response code.
In this scenario, the x-auth-status-code header is appended to the request by the ext authz HTTP filter with the value 200. 
The upstream response status code is 200.  
So there is a match. No metric is emitted. 

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
No metric is emitted. 


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


