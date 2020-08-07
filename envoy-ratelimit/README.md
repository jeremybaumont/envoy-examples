# envoy-ratelimit

## What

This is a simple sandbox with envoy as front proxy and [httpbin](https://eu.httpbin.org/) as upstream 
service, and a rate limit service mock (that for now always allows the traffic). `

## Why

You can easily test and reproduce rate limit configurations.

## How

* Start the local world:
```bash
docker-compose up --build
```

* In another terminal, make an http request to test your whatever envoy config with a http client like [HTTPie](https://httpie.org/) or curl.
```bash
http localhost:8080/anything
```

